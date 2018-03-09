package main

import (
	"go/importer"
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

type GetName interface {
	GetName() string
}

type GetDescription interface {
	GetDescription() string
}

type GetPath interface {
	GetPath() string
}

type Get interface {
	Get(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
}

type Put interface {
	Put(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
}

func CheckUnimplimented(obj plugin.Symbol) ([]string, string) {
	var unimplemented []string
	_, ok := obj.(GetName)
	if !ok {
		return nil, "getName not Implemented"
	}
	_, ok = obj.(GetDescription)
	if !ok {
		return nil, "getDescription not Implemented"
	}
	_, ok = obj.(GetPath)
	if !ok {
		return nil, "getPath not Implemented"
	}
	_, ok = obj.(Get)
	if !ok {
		unimplemented = append(unimplemented, "get")
	}
	_, ok = obj.(Put)
	if !ok {
		unimplemented = append(unimplemented, "put")
	}
	return unimplemented, ""
}

func hasCrud(inter Crud, ctx config.AppContext, router *mux.Router) {
	var handler http.Handler
	println()
	color.Blue("Attaching Endpoints from plugins: ")
	handler = MakeHandler(ctx, inter.Get)
	router.
		Methods("GET").
		Path(inter.GetPath()).
		Name(inter.GetName()).
		Handler(handler)
	color.Green("%s -- %s -- %s : %s", inter.GetName(), "GET", inter.GetPath(), inter.GetDescription())
}

func noCrud(name string, obj plugin.Symbol) {
	println()
	color.Red("%s does not implement the Crud interface", name)
	m, err := CheckUnimplimented(obj)
	if err != "" {
		println(err)
	}
	color.Yellow("%s is missing the following methods:", name)
	for _, method := range m {
		color.Yellow(method)
	}
	println()
}

func ApplyPlugin(plug string, router *mux.Router, ctx config.AppContext) {

	prefix := "_plugins/out/"
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	importString := dir[len(os.Getenv("GOPATH")+"/src/"):] + "/_plugins/main"
	p, err := plugin.Open(prefix + plug + ".so")
	pkg, err := importer.For("source", nil).Import(importString)
	if err != nil {
		panic(err)
	}

	// Get Package Definition
	for _, name := range pkg.Scope().Names() {
		definition := pkg.Scope().Lookup(name)
		// Check for structure
		prefix := strings.HasPrefix(definition.Type().Underlying().String(), "struct")
		exported := definition.Exported()
		if prefix && exported {
			obj, err := p.Lookup(name)
			if err != nil {
				continue
			}
			switch inter := obj.(type) {
			case Crud:
				hasCrud(inter, ctx, router)
			default:
				noCrud(name, obj)
			}
		}
	}
	println()
}
