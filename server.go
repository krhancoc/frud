// Package frud
package main

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/krhancoc/frud/config"
	"github.com/krhancoc/frud/database"
	"github.com/krhancoc/frud/middleware"
	"github.com/krhancoc/frud/plug"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	"github.com/unrolled/secure"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// StartServer Wraps the mux Router and uses the Negroni Middleware
func StartServer(path string) *http.Server {
	//Load up Logger
	// Load up database
	conf, err := config.LoadConfig(path)
	if err != nil {
		panic(err)
	}
	log.Debug("Setting up database")
	db, err := database.CreateDatabase(conf)
	if err != nil {
		panic(err)
	}

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

	log.Info("Creating Plugin Manager")
	plugManager, err := plug.CreatePlugManager(conf.Manager)
	if err != nil {
		panic(err)
	}
	log.Info("Attaching Routes")
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

	ln, err := net.Listen("tcp", ":"+ctx.Port)
	if err != nil {
		panic(err)
	}
	srv := http.Server{
		Addr:    ":" + ctx.Port,
		Handler: n,
	}

	go srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
	_, err = http.Get("http://localhost:" + ctx.Port)
	for err == nil {
		return &srv
	}
	srv.Shutdown(nil)
	panic(err)
	return nil
}
