package plug

import (
	"fmt"
	"go/importer"
	"go/types"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/krhancoc/frud/config"
)

func getImportString(conf *config.PlugConfig) string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir[len(os.Getenv("GOPATH")+"/src/"):] + "/" + conf.PathToCode
}

func plugPack(plug string, conf *config.PlugConfig) (*plugin.Plugin, *types.Package, error) {

	importString := getImportString(conf) + plug
	p, err := plugin.Open(conf.PathToCompiled + plug + ".so")
	if err != nil {
		return nil, nil, err
	}
	pkg, err := importer.For("source", nil).Import(importString)
	if err != nil {
		return nil, nil, err
	}
	return p, pkg, nil
}

func isPlugin(o types.Object) bool {

	prefix := strings.HasPrefix(o.Type().Underlying().String(), "struct")
	exported := o.Exported()
	return prefix && exported
}

// CreatePlugManager creates the manager object for Plugins
func CreatePlugManager(conf *config.ManagerConfig) (*Manager, error) {

	var plugManager *Manager
	var plugs []*Plug

	for _, plug := range conf.Plugs.Names {

		var unimplemented []string

		thisPlug := CreatePlug()
		modelFound := false
		handlerFound := false

		p, pkg, err := plugPack(plug, conf.Plugs)
		if err != nil {
			return plugManager, err
		}

		missing := thisPlug.SetModel(p)
		if len(missing) == 0 {
			modelFound = true
		}
		// Check for Handlers
		for _, name := range pkg.Scope().Names() {
			definition := pkg.Scope().Lookup(name)
			// Check for structure
			if isPlugin(definition) {
				obj, err := p.Lookup(name)
				if err != nil {
					continue
				}
				err = thisPlug.SetDefinition(name, obj)
				if err != nil {
					return plugManager, err
				}
				unimplemented = thisPlug.SetCrud(name, obj)
				if len(unimplemented) == 0 {
					handlerFound = true
				}
				plugs = append(plugs, &thisPlug)
			}
		}

		if thisPlug.Name == "" {
			thisPlug.SetDefaultDefinition(plug)
		}

		if !modelFound && !handlerFound {
			return plugManager, fmt.Errorf(`
				%s: Requires a Model implementation or CRUD implementation.
				Missing Model Functions: %s
				`, plug, strings.Join(missing, ","))
		} else if handlerFound {
			color.Yellow("Plugin found using Handler Method - %s: %s", thisPlug.Name, thisPlug.Description)
		} else if modelFound {
			color.Yellow("Plugin found using Model Method - %s: %s", thisPlug.Name, thisPlug.Description)
		}

	}
	println()
	plugManager = &Manager{
		Plugs: plugs,
	}
	return plugManager, nil
}

// AttachRoutes to your router!!
func (manager *Manager) AttachRoutes(router *mux.Router, ctx config.AppContext) error {

	color.Cyan("Attaching routes...")
	println()
	for _, plug := range manager.Plugs {
		color.Yellow("Plugin %s: %s", plug.Name, plug.Description)
		color.Yellow("---------------------")
		inter := *plug.Crud
		if inter != nil {
			methods := map[string]HandlerFunc{
				"Get":    inter.Get,
				"Post":   inter.Post,
				"Put":    inter.Put,
				"Delete": inter.Delete,
			}
			for method, f := range methods {
				var handler http.Handler
				handler = MakeHandler(ctx, f)
				router.
					Methods(method).
					Path(plug.EntryPoint).
					Name(plug.Name).
					Handler(handler)
				color.Green("%s -- %s -- %s", plug.Name, method, plug.EntryPoint)
			}
			println("\n")
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
