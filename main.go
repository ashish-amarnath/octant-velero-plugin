package main

import (
	"log"

	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"

	vplugin "github.com/ashish-amarnath/octant-velero-plugin/internal/plugin"
)

func main() {
	log.SetPrefix("")

	vp := &vplugin.VeleroPlugin{}

	capabilities := &plugin.Capabilities{
		IsModule: true,
		ActionNames: []string{
			vplugin.NamespaceAction,
			vplugin.DeleteAction,
		},
	}

	options := []service.PluginOption{
		service.WithNavigation(vp.Navigation, vp.InitRoutes),
		service.WithActionHandler(vp.ActionHandler),
	}

	p, err := service.Register(vplugin.Name, vplugin.Description, capabilities, options...)
	if err != nil {
		log.Fatalf("Failed to register %s: %v", vplugin.Name, err)
	}

	log.Printf("Starting plugin %s", vplugin.Name)
	p.Serve()
}
