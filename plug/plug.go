package plug

import (
	"reflect"
)

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
	t := reflect.TypeOf((*Definition)(nil)).Elem()
	for i := 0; i < t.NumMethod(); i++ {
		f := t.Method(i).Name
		check := reflect.TypeOf(o)
		_, ok := check.MethodByName(f)
		if !ok {
			unimplemented = append(unimplemented, f)
		}
	}
	return unimplemented
}

func (plug *Plug) CheckUnimplimented(i interface{}) []string {

	var unimplemented []string
	t := reflect.TypeOf(i).Elem()
	for i := 0; i < t.NumMethod(); i++ {
		f := t.Method(i).Name
		check := reflect.TypeOf(*plug.Main)
		_, ok := check.MethodByName(f)
		if !ok {
			unimplemented = append(unimplemented, f)
		}
	}
	return unimplemented
}
