package main

import (
	"net/http"
	"plugin"

	"github.com/krhancoc/frud/config"
)

type GetName interface {
	GetName() string
}

type GetDescription interface {
	GetDescription() string
}

type GetPath interface {
	GetPath() string
}

type Get interface {
	Get(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
}

type Put interface {
	Put(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
}

func CheckUnimplimented(obj plugin.Symbol) ([]string, string) {
	var unimplemented []string
	_, ok := obj.(GetName)
	if !ok {
		return nil, "getName not Implemented"
	}
	_, ok = obj.(GetDescription)
	if !ok {
		return nil, "getDescription not Implemented"
	}
	_, ok = obj.(GetPath)
	if !ok {
		return nil, "getPath not Implemented"
	}
	_, ok = obj.(Get)
	if !ok {
		unimplemented = append(unimplemented, "get")
	}
	_, ok = obj.(Put)
	if !ok {
		unimplemented = append(unimplemented, "put")
	}
	return unimplemented, ""
}
