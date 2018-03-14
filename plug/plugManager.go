package plug

import (
	"net/http"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/krhancoc/frud/config"
)

// CreatePlugManager creates the manager object for Plugins
func CreatePlugManager(conf *config.ManagerConfig) (*Manager, error) {

	var plugManager *Manager
	var plugs []*Plug

	for _, plug := range conf.Plugs {

		p, err := createPlug(plug)
		if err != nil {
			panic(err)
		}
		plugs = append(plugs, p)
	}
	println()
	plugManager = &Manager{
		Plugs: plugs,
	}
	return plugManager, nil
}

// AttachRoutes to your router!!
func (manager *Manager) AttachRoutes(router *mux.Router, ctx config.AppContext) error {

	println()
	for _, plug := range manager.Plugs {
		color.Yellow("Plugin %s: %s", plug.Name, plug.Description)
		color.Yellow("---------------------")
		var methods map[string]http.HandlerFunc
		if plug.Crud != nil {
			inter := *plug.Crud
			methods = map[string]http.HandlerFunc{
				"Get":    MakeHandler(ctx, inter.Get),
				"Post":   MakeHandler(ctx, inter.Post),
				"Put":    MakeHandler(ctx, inter.Put),
				"Delete": MakeHandler(ctx, inter.Delete),
			}
		} else {
			methods = map[string]http.HandlerFunc{
				"Get":    makeGenericHandler(ctx, *plug, get),
				"Post":   makeGenericHandler(ctx, *plug, post),
				"Delete": makeGenericHandler(ctx, *plug, delete),
				"Put":    makeGenericHandler(ctx, *plug, put),
			}
		}
		for method, f := range methods {
			router.
				Methods(method).
				Path(plug.Path).
				Name(plug.Name).
				Handler(f)
			color.Green("%s -- %s -- %s", plug.Name, method, plug.Path)
		}
	}
	return nil
}

func (*Manager) InitGenericRoutes(router *mux.Router, conf *config.ManagerConfig, ctx config.AppContext) error {

	for _, generic := range conf.Generics {
		switch generic {
		case "healthcheck":
			router.
				Methods("Get").
				Path("/health").
				Name("Health Check").
				Handler(MakeHandler(ctx, HealthCheck))
		case "login":
			println(generic)
		}
	}
	return nil
}
