package main

import (
	"net/http"

	"github.com/krhancoc/frud/config"
)

type EntryEndpoint struct {
	Data *config.Endpoint
}

var EntryEndpointObject EntryEndpoint = EntryEndpoint{
	Data: &config.Endpoint{
		Name:        "Entry",
		Description: "Entry object description",
		Path:        "/entry",
	},
}

func (*EntryEndpoint) Get(w http.ResponseWriter, req *http.Request, ctx config.AppContext) {
	ctx.Render.Text(w, 200, "HELLO")
}

func (*EntryEndpoint) Post(w http.ResponseWriter, req *http.Request, ctx config.AppContext) {
	ctx.Render.Text(w, 200, "HELLO")
}

func (*EntryEndpoint) Delete(w http.ResponseWriter, req *http.Request, ctx config.AppContext) {
	ctx.Render.Text(w, 200, "HELLO")
}

func (*EntryEndpoint) Put(w http.ResponseWriter, req *http.Request, ctx config.AppContext) {
	ctx.Render.Text(w, 200, "HELLO")
}
