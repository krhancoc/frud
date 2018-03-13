package config

import (
	"github.com/unrolled/render"
)

type Driver interface {
	Connect() interface{}
	MakeRequest(string, map[string]string, []*Field) error
}

// Configuration object that acts as the parent to others
type Configuration struct {
	Context  *Context       `json:"context"`
	Database *Database      `json:"database"`
	Manager  *ManagerConfig `json:"manager"`
}

type ManagerConfig struct {
	Generics []string      `json:"generics"`
	Plugs    []*PlugConfig `json:"plugins"`
}

type PlugConfig struct {
	PathToCode     string   `json:"pathtocode,omitempty"`
	PathToCompiled string   `json:"pathtocompiled,omitempty"`
	Name           string   `json:"name"`
	Description    string   `json:"description,omitempty"`
	Path           string   `json:"path,omitempty"`
	Model          []*Field `json:"model,omitempty"`
}

type Field struct {
	Key       string `json:"key"`
	ValueType string `json:"value_type"`
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

type Endpoint struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
}
