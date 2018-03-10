package config

import (
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/unrolled/render"
)

// Configuration object that acts as the parent to others
type Configuration struct {
	Context  *Context    `json:"context"`
	Database *Database   `json:"database"`
	Plugins  *PlugConfig `json:"plugins"`
}

type PlugConfig struct {
	PathToCode     string   `json:"pathtocode,omitempty"`
	PathToCompiled string   `json:"pathtocompiled,omitempty"`
	Names          []string `json:"names"`
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
	Connection *bolt.Conn
	Render     *render.Render
	Version    string
	Port       string
}

type Endpoint struct {
	Name string `json:"name"`
}
