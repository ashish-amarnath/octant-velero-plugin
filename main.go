package main

import (
	"fmt"
	"log"

	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
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
)

func main() {
	capabilities := &plugin.Capabilities{
		IsModule: true,
	}

	// Set up what should happen when Octant calls this plugin.
	options := []service.PluginOption{
		service.WithNavigation(veleroPluginNavigation, initPluginRoutes),
	}

	p, err := service.Register(pluginName, pluginDesc, capabilities, options...)
	if err != nil {
		log.Fatalf("Failed to register %s: %v", pluginName, err)
	}

	log.Printf("Starting plugin %s", pluginName)
	p.Serve()
}

func veleroPluginNavigation(request *service.NavigationRequest) (navigation.Navigation, error) {
	return navigation.Navigation{
		Title:    "Velero Dashboard",
		Path:     request.GeneratePath("velero-dashboard"),
		Children: []navigation.Navigation{},
		IconName: "cloud",
	}, nil
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

func itemToRow(iType string) component.TableRow {
	row := component.TableRow{}
	cols := getColumnsForTable(iType)
	for _, c := range cols {
		row[c] = component.NewText(c)
	}
	return row
}

func initPluginRoutes(router *service.Router) {
	gen := func(name, accessor, requestPath string) component.Component {
		tableName := fmt.Sprintf("Velero %s", name)
		placeholder := fmt.Sprintf("We could not find any %s!", tableName)

		table := component.NewTable(tableName, placeholder,
			component.NewTableCols(getColumnsForTable(accessor)...))
		// TODO:
		// 1. get Velero <accessor> as a list
		// 2. for each item in list, itemToRow
		table.Add(itemToRow(accessor))

		return table
	}

	router.HandleFunc("*", func(request service.Request) (component.ContentResponse, error) {
		components := []component.Component{
			gen("Backups", "backup", request.Path()),
			gen("Restores", "restore", request.Path()),
			gen("Schedules", "schedule", request.Path()),
		}

		contentResponse := component.NewContentResponse(component.TitleFromString("Velero Dashboard"))
		contentResponse.Add(components...)

		return *contentResponse, nil
	})

}
