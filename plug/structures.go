package plug

import (
	"go/types"
	"net/http"

	"github.com/krhancoc/frud/config"
)

type Model struct {
	Data interface{}
}

// Manager object to hold our plugins
type Manager struct {
	Plugs []*Plug
}

// Plug object to hold each plugin
type Plug struct {
	Name        string
	Description string
	EntryPoint  string
	Package     *types.Package
	Crud        *Crud
	Model       map[string](func(map[string]string) interface{})
}

// Crud interface for the functions required by our API objects
type Crud interface {
	Get(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
	Put(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
	Delete(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
	Post(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
}

// Definition interface for plugins so they can describe themselves
type definition interface {
	GetName() string
	GetPath() string
	GetDescription() string
}

// Endpoint interface is the wrapper for functionality and definition.  Splitting them up
// for now as I will want to add Default handlers for Get, Put, Delete, Post functions in
// case users do not want to make their own.
type endpoint interface {
	Crud
	definition
}

type modelFunctions interface {
	Get(i interface{}) string
	Put(i interface{}) string
	Delete(i interface{}) string
	Post(i interface{}) string
}

// HandlerFunc type for the unwrapped version of a handler function.
type HandlerFunc func(http.ResponseWriter, *http.Request, config.AppContext)

// MakeHandler will wrap a Application Context up into a handler so that users
// have access to things like their database connection etc in each handler.
func MakeHandler(ctx config.AppContext, fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, ctx)
	}
}
