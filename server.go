package main

import (
	"log"
	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/krhancoc/frud/config"
	"github.com/krhancoc/frud/database"
	"github.com/krhancoc/frud/middleware"
	"github.com/krhancoc/frud/plug"
	"github.com/unrolled/render"
	"github.com/unrolled/secure"
)

// Crud endpoint

// StartServer Wraps the mux Router and uses the Negroni Middleware
func StartServer(path string) {

	// Load up database
	conf := config.LoadConfig(path)
	db := database.CreateDatabase(conf.Database)

	// Create App context
	ctx := config.AppContext{
		Driver:  db.(config.Driver),
		Render:  render.New(),
		Version: conf.Context.Version,
		Port:    strconv.Itoa(conf.Context.Port),
	}

	//Instantiate middleware
	router := mux.NewRouter().StrictSlash(true)
	router.Use(mux.MiddlewareFunc(middleware.Converter))

	plugManager, err := plug.CreatePlugManager(conf.Manager)
	if err != nil {
		panic(err)
	}
	err = plugManager.AttachRoutes(router, ctx)
	if err != nil {
		panic(err)
	}
	err = plugManager.InitGenericRoutes(router, conf.Manager, ctx)
	if err != nil {
		panic(err)
	}
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
