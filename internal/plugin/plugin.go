package plugin

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

const (
	Name        = "velero.io/octant-plugin"
	Description = "Velero Dashboard Plugin"
)

var (
	NamespaceAction = action.RequestSetNamespace
	DeleteAction    = "velero.io/backupDelete"
)

type VeleroPlugin struct {
	currentNamespace string
	rmu              sync.RWMutex
}

func (v *VeleroPlugin) SetCurrentNamespace(namespace string) {
	v.rmu.Lock()
	v.rmu.Unlock()
	v.currentNamespace = namespace
}

func (v *VeleroPlugin) CurrentNamespace() string {
	v.rmu.RLock()
	v.rmu.RUnlock()
	return v.currentNamespace
}

func (v *VeleroPlugin) ActionHandler(request *service.ActionRequest) error {
	switch request.ActionName {
	case NamespaceAction:
		currentNamespace, err := request.Payload.String("namespace")
		if err != nil {
			return err
		}
		v.SetCurrentNamespace(currentNamespace)
		return nil
	case DeleteAction:
		log.Printf("handle delete action here")
		return nil
	default:
		return fmt.Errorf("no action %s registered for plugin %s", request.ActionName, Name)
	}
}

func (v *VeleroPlugin) Navigation(request *service.NavigationRequest) (navigation.Navigation, error) {
	return navigation.Navigation{
		Title:    "Velero Dashboard",
		Path:     request.GeneratePath("velero-dashboard"),
		Children: []navigation.Navigation{},
		IconName: "cloud",
	}, nil
}

func (v *VeleroPlugin) InitRoutes(router *service.Router) {
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
