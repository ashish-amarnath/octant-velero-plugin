package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"

	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

const (
	pluginName = "velero.io/octant-plugin"
	pluginDesc = "Velero Dashboard Plugin"
)

var (
	backupColumns   = []string{"Name", "Status", "Created", "Expires", "Storage Location", "Selector"}
	restoreColumns  = []string{"Name", "Backup", "Status", "Warnings", "Errors", "Created", "Selector"}
	scheduleColumns = []string{"Name", "Status", "Created", "Schedule", "Backup TTL", "Last Backup", "Selector"}

	backupGvk   = schema.GroupVersionKind{Group: "velero.io", Version: "v1", Kind: "Backup"}
	restoreGvk  = schema.GroupVersionKind{Group: "velero.io", Version: "v1", Kind: "Restore"}
	scheduleGvk = schema.GroupVersionKind{Group: "velero.io", Version: "v1", Kind: "Schedule"}

	deleteAction = "velero.io/backupDelete"
)

type veleroPlugin struct {
	currentNamespace string
	rmu              sync.RWMutex
}

func main() {
	log.SetPrefix("")

	vp := &veleroPlugin{}

	capabilities := &plugin.Capabilities{
		IsModule:    true,
		ActionNames: []string{action.RequestSetNamespace, deleteAction},
	}

	// Set up what should happen when Octant calls this plugin.
	options := []service.PluginOption{
		service.WithNavigation(vp.Navigation, vp.InitRoutes),
		service.WithActionHandler(vp.SetNamespaceHandler),
	}

	p, err := service.Register(pluginName, pluginDesc, capabilities, options...)
	if err != nil {
		log.Fatalf("Failed to register %s: %v", pluginName, err)
	}

	log.Printf("Starting plugin %s", pluginName)
	p.Serve()
}

func (v *veleroPlugin) CurrentNamespace() string {
	v.rmu.RLock()
	defer v.rmu.RUnlock()
	return v.currentNamespace
}

func (v *veleroPlugin) SetNamespaceHandler(request *service.ActionRequest) error {
	switch request.ActionName {
	case action.RequestSetNamespace:
		currentNamespace, err := request.Payload.String("namespace")
		if err != nil {
			return err
		}
		v.rmu.Lock()
		defer v.rmu.Unlock()
		v.currentNamespace = currentNamespace
		return nil
	default:
		return fmt.Errorf("no action %s registered for plugin %s", request.ActionName, pluginName)
	}
}

func (v *veleroPlugin) Navigation(request *service.NavigationRequest) (navigation.Navigation, error) {
	return navigation.Navigation{
		Title:    "Velero Dashboard",
		Path:     request.GeneratePath("velero-dashboard"),
		Children: []navigation.Navigation{},
		IconName: "cloud",
	}, nil
}

func (v *veleroPlugin) InitRoutes(router *service.Router) {
	gen := func(ctx context.Context, client service.Dashboard, name, accessor, requestPath string) component.Component {
		tableName := fmt.Sprintf("Velero %s", name)
		placeholder := fmt.Sprintf("We could not find any %s!", tableName)

		table := component.NewTable(tableName, placeholder,
			component.NewTableCols(getColumnsForTable(accessor)...))

		gvk, err := getGVKForAccessor(accessor)
		if err != nil {
			// Handle bad GVK
			log.Printf("gvk failed %s", accessor)
			return nil
		}

		key := store.KeyFromGroupVersionKind(gvk)
		key.Namespace = v.CurrentNamespace()

		objects, err := client.List(ctx, key)
		if err != nil {
			log.Printf("listing failed for %+v, %s", key, err)
			return table
		}

		for _, object := range objects.Items {
			table.Add(itemToRow(accessor, object))
		}

		return table
	}

	router.HandleFunc("*", func(request service.Request) (component.ContentResponse, error) {
		components := []component.Component{
			gen(request.Context(), request.DashboardClient(), "Backups", "backup", request.Path()),
			gen(request.Context(), request.DashboardClient(), "Restores", "restore", request.Path()),
			gen(request.Context(), request.DashboardClient(), "Schedules", "schedule", request.Path()),
		}

		contentResponse := component.NewContentResponse(component.TitleFromString("Velero Dashboard"))
		contentResponse.Add(components...)

		return *contentResponse, nil
	})

}

func getColumnsForTable(name string) []string {
	switch name {
	case "backup":
		return backupColumns
	case "restore":
		return restoreColumns
	case "schedule":
		return scheduleColumns
	default:
		return []string{}
	}
}

func getGVKForAccessor(accessor string) (schema.GroupVersionKind, error) {
	switch accessor {
	case "backup":
		return backupGvk, nil
	case "restore":
		return restoreGvk, nil
	case "schedule":
		return scheduleGvk, nil
	default:
		return schema.GroupVersionKind{}, fmt.Errorf("bad accessor, no GVK found")
	}
}

type rowPrinter func(unstructured.Unstructured) (component.TableRow, error)

func getPrinterForAccessor(accessor string) (rowPrinter, error) {
	switch accessor {
	case "backup":
		return backupPrinter, nil
	case "restore":
		return restorePrinter, nil
	case "schedule":
		return schedulePrinter, nil
	default:
		return nil, fmt.Errorf("bad accessor, no printer found")
	}
}

func itemToRow(accessor string, object unstructured.Unstructured) component.TableRow {
	printer, err := getPrinterForAccessor(accessor)
	if err != nil {
		log.Printf("loading printer error: %s", err)
		return nil
	}
	row, err := printer(object)
	if err != nil {
		log.Printf("printing error: %s", err)
	}
	return row
}

func restorePrinter(object unstructured.Unstructured) (component.TableRow, error) {
	restore := &velerov1.Restore{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, restore)
	if err != nil {
		return component.TableRow{}, err
	}
	row := component.TableRow{}
	//TODO: Properly format these entries
	row["Name"] = component.NewText(fmt.Sprintf("%s", restore.Name))
	row["Backup"] = component.NewText(fmt.Sprintf("%s", restore.Spec.BackupName))
	row["Status"] = component.NewText(fmt.Sprintf("%s", restore.Status.Phase))
	row["Warnings"] = component.NewText(fmt.Sprintf("%d", restore.Status.Warnings))
	row["Errors"] = component.NewText(fmt.Sprintf("%d", restore.Status.Errors))
	row["Created"] = component.NewText(fmt.Sprintf("%s", restore.CreationTimestamp))
	row["Selector"] = component.NewLabels(restore.Spec.LabelSelector.MatchLabels)
	return row, nil
}

func schedulePrinter(object unstructured.Unstructured) (component.TableRow, error) {
	schedule := &velerov1.Schedule{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, schedule)
	if err != nil {
		return component.TableRow{}, err
	}
	row := component.TableRow{}
	//TODO: Properly format these entries
	row["Name"] = component.NewText(fmt.Sprintf("%s", schedule.Name))
	row["Status"] = component.NewText(fmt.Sprintf("%s", schedule.Status.Phase))
	row["Created"] = component.NewText(fmt.Sprintf("%s", schedule.CreationTimestamp))
	row["Backup TTL"] = component.NewText(fmt.Sprintf("%s", schedule.Spec.Template.TTL))
	row["Last Backup"] = component.NewText(fmt.Sprintf("%s", schedule.Status.LastBackup))
	row["Selector"] = component.NewLabels(schedule.Spec.Template.LabelSelector.MatchLabels)
	return row, nil
}

func backupPrinter(object unstructured.Unstructured) (component.TableRow, error) {
	backup := &velerov1.Backup{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, backup)
	if err != nil {
		return component.TableRow{}, err
	}
	backupColumns = []string{"Name", "Status", "Created", "Expires", "Storage Location", "Selector"}

	row := component.TableRow{}
	//TODO: Properly format these entries
	row["Name"] = component.NewText(fmt.Sprintf("%s", backup.Name))
	row["Status"] = component.NewText(fmt.Sprintf("%s", backup.Status.Phase))
	row["Created"] = component.NewText(fmt.Sprintf("%s", backup.CreationTimestamp))
	row["Expires"] = component.NewText(fmt.Sprintf("%s", backup.Spec.TTL))
	row["Storage Location"] = component.NewText(fmt.Sprintf("%s", backup.Spec.StorageLocation))
	row["Selector"] = component.NewLabels(backup.Spec.LabelSelector.MatchLabels)
	return row, nil
}
