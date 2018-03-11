package main

import (
	"net/http"

	"github.com/krhancoc/frud/config"
)

type EntryModel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type EntryEndpoint struct {
	Name        string
	Description string
	Path        string
}

var EntryEndpointObject EntryEndpoint = EntryEndpoint{
	Name:        "Entry",
	Description: "Entry object description",
	Path:        "/entry",
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

func (e *EntryEndpoint) GetName() string {
	return e.Name
}

func (e *EntryEndpoint) GetDescription() string {
	return e.Description
}

func (e *EntryEndpoint) GetPath() string {
	return e.Path
}
