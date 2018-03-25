package config

import (
	"github.com/unrolled/render"
)

var allowedGenerics = []string{
	"healthcheck",
}

var allowedTypes = []string{
	"int",
	"string",
}

// Driver interface is what all database drivers must implement, it allows the server
// to interact with the database without actually caring what the database actually is.
type Driver interface {
	MakeRequest(*DBRequest) (interface{}, error)
	ConvertToDriverError(error) error
}

// Context of the server itself
type Context struct {
	Port    int    `json:"port"`
	Version string `json:"version"`
}

// Database configuration
type Database struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

// AppContext is the structure that holds the current context of the app
// will be passed in the handlers
type AppContext struct {
	Driver  Driver
	Render  *render.Render
	Version string
	Port    string
}
