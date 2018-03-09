package main

import (
	"go/importer"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plugin"
	"strconv"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/krhancoc/frud/config"
	"github.com/krhancoc/frud/database"
	"github.com/krhancoc/frud/middleware"
	"github.com/unrolled/render"
	"github.com/unrolled/secure"
)

// Crud endpoint
type Crud interface {
	Get(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
	// Put(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
	// Delete(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
	// Post(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
	GetName() string
	GetDescription() string
	GetPath() string
}

type HandlerFunc func(http.ResponseWriter, *http.Request, config.AppContext)

func MakeHandler(ctx config.AppContext, fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, ctx)
	}
}

// StartServer Wraps the mux Router and uses the Negroni Middleware
func StartServer(path string) {

	// Load up database
	conf := config.LoadConfig(path)
	connect := database.Connect(conf.Database)
	defer connect.Close()

	// Create App context
	ctx := config.AppContext{
		Connection: &connect,
		Render:     render.New(),
		Version:    conf.Context.Version,
		Port:       strconv.Itoa(conf.Context.Port),
	}

	//Instantiate middleware
	router := mux.NewRouter().StrictSlash(true)
	router.Use(mux.MiddlewareFunc(middleware.Converter))

	prefix := "_plugins/out/"
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	// Import the package
	importString := dir[len(os.Getenv("GOPATH")+"/src/"):] + "/_plugins/main"
	for _, plug := range conf.Plugins {
		p, err := plugin.Open(prefix + plug + ".so")
		pkg, err := importer.For("source", nil).Import(importString)
		if err != nil {
			panic(err)
		}
		for _, name := range pkg.Scope().Names() {
			definition := pkg.Scope().Lookup(name)
			if strings.HasPrefix(definition.Type().Underlying().String(), "struct") && definition.Exported() {
				obj, err := p.Lookup(name)
				if err != nil {
					continue
				}
				switch inter := obj.(type) {
				case Crud:
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
				default:
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
			}
		}
		println()

	}
	// security
	isDevelopment := true

	secureMiddleware := secure.New(secure.Options{
		IsDevelopment:      isDevelopment, // This will cause the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains options to be ignored during development. When deploying to production, be sure to set this to false.
		AllowedHosts:       []string{},    // AllowedHosts is a list of fully qualified domain names that are allowed (CORS)
		ContentTypeNosniff: true,          // If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
		BrowserXssFilter:   true,          // If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
	})

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.UseHandler(router)
	log.Println("===> Starting app (v" + ctx.Version + ") on port " + ctx.Port)
	n.Run(":" + ctx.Port)
}
