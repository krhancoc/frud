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

type HandlerFunc func(http.ResponseWriter, *http.Request, config.AppContext)

func MakeHandler(ctx config.AppContext, fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, ctx)
	}
}

// Crud interface for the functions required by our API objects
type Crud interface {
	Get(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
	Put(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
	Delete(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
	Post(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
}

// PlugManager object to hold our plugins
type Manager struct {
	Plugs []*Plug
}

// Plug object to hold each plugin
type Plug struct {
	Name        string
	Description string
	EntryPoint  string
	Package     *types.Package
	Main        *plugin.Symbol
}

type Definition interface {
	GetName() string
	GetPath() string
	GetDescription() string
}

func getImportString(conf *config.PlugConfig) string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir[len(os.Getenv("GOPATH")+"/src/"):] + "/" + conf.PathToCode
}

func plugPack(plug string, conf *config.PlugConfig) (*plugin.Plugin, *types.Package, error) {

	importString := getImportString(conf)
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

		p, pkg, err := plugPack(plug, conf.Plugs)
		if err != nil {
			return plugManager, err
		}

		// Get Package Definition
		for _, name := range pkg.Scope().Names() {
			definition := pkg.Scope().Lookup(name)
			// Check for structure
			if isPlugin(definition) {
				obj, err := p.Lookup(name)
				if err != nil {
					continue
				}
				unimplimented := CheckDefinition(obj)
				if len(unimplimented) > 0 {
					return plugManager, fmt.Errorf("Unimplemented definition functions in - %s, %s", name, strings.Join(unimplimented, ","))
				}
				inter := obj.(Definition)
				thisPlug := Plug{
					Name:        inter.GetName(),
					Description: inter.GetDescription(),
					EntryPoint:  inter.GetPath(),
					Package:     pkg,
					Main:        &obj,
				}
				color.Yellow("Plugin found - %s: %s", thisPlug.Name, thisPlug.Description)
				plugs = append(plugs, &thisPlug)
			}
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
		unimplimented := plug.CheckUnimplimented((*Crud)(nil))
		if len(unimplimented) > 0 {
			return fmt.Errorf("Plug: %s, unimplimented - %s", plug.Name, strings.Join(unimplimented, ","))
		}
		inter := (*plug.Main).(Crud)
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
	return nil
}

func (*Manager) InitGenericRoutes(router *mux.Router, conf *config.ManagerConfig) error {

	for _, generic := range conf.Generics {
		switch generic {
		case "healthcheck":
			println(generic)
		case "login":
			println(generic)
		}
	}
	return nil
}
