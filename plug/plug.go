package plug

import (
	"net/http"

	"github.com/krhancoc/frud/config"
)

type get interface {
	Get(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
}

type put interface {
	Put(w http.ResponseWriter, req *http.Request, ctx config.AppContext)
}

type getName interface {
	GetName() string
}

type getDescription interface {
	GetDescription() string
}

type getPath interface {
	GetPath() string
}

func CheckDefinition(o interface{}) []string {
	var unimplemented []string
	_, ok := o.(getName)
	if !ok {
		unimplemented = append(unimplemented, "GetName")
	}
	_, ok = o.(getDescription)
	if !ok {
		unimplemented = append(unimplemented, "GetDescription")
	}
	_, ok = o.(getName)
	if !ok {
		unimplemented = append(unimplemented, "GetName")
	}
	return unimplemented
}

func (plug *Plug) CheckUnimplimented() []string {
	var unimplemented []string
	_, ok := (*plug.Main).(get)
	if !ok {
		unimplemented = append(unimplemented, "get")
	}
	// _, ok = (*plug.Main).(put)
	// if !ok {
	// 	unimplemented = append(unimplemented, "put")
	// }
	return unimplemented
}
