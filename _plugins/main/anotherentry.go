package main

import (
	"net/http"

	"github.com/krhancoc/frud/config"
)

type AnotherEntryModel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type AnotherEntryEndpoint struct {
	Name        string
	Description string
	Path        string
}

var AnotherEntryEndpointObject AnotherEntryEndpoint = AnotherEntryEndpoint{
	Name:        "AnotherEntry",
	Description: "AnotherEntry object description",
	Path:        "/another",
}

func (*AnotherEntryEndpoint) Get(w http.ResponseWriter, req *http.Request, ctx config.AppContext) {
	ctx.Render.Text(w, 200, "ANOTHER HELLO")
}

func (*AnotherEntryEndpoint) Post(w http.ResponseWriter, req *http.Request, ctx config.AppContext) {
	ctx.Render.Text(w, 200, "ANOTHER HELLO")
}

func (*AnotherEntryEndpoint) Put(w http.ResponseWriter, req *http.Request, ctx config.AppContext) {
	ctx.Render.Text(w, 200, "ANOTHER HELLO")
}

func (*AnotherEntryEndpoint) Delete(w http.ResponseWriter, req *http.Request, ctx config.AppContext) {
	ctx.Render.Text(w, 200, "ANOTHER HELLO")
}

func (e *AnotherEntryEndpoint) GetName() string {
	return e.Name
}

func (e *AnotherEntryEndpoint) GetDescription() string {
	return e.Description
}

func (e *AnotherEntryEndpoint) GetPath() string {
	return e.Path
}
