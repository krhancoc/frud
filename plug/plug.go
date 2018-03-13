package plug

import (
	"fmt"
	"plugin"
	"reflect"
	"strings"
)

func createPlug(name string) Plug {
	return Plug{
		Name:        name,
		Description: "Default description " + name,
		EntryPoint:  "/" + name,
		Crud:        nil,
		Model:       nil,
	}
}

// CheckUnimplimented will check if interface obj has all the functions asked for by i, it will then
// output a list of the functions not implimented.
// Please note that obj has to be a pointer to the struct in question and not the struct object itself.
func checkUnimplimented(obj interface{}, i interface{}) []string {
	var unimplemented []string
	t := reflect.TypeOf(i).Elem()
	for i := 0; i < t.NumMethod(); i++ {
		f := t.Method(i).Name
		check := reflect.TypeOf(obj)
		_, ok := check.MethodByName(f)
		if !ok {
			unimplemented = append(unimplemented, f)
		}
	}
	return unimplemented
}

func (w *Plug) setDefinition(name string, obj interface{}) error {
	unimplimented := checkUnimplimented(obj, (*definition)(nil))
	if len(unimplimented) > 0 {
		return fmt.Errorf("Unimplemented definition functions in - %s, %s", name, strings.Join(unimplimented, ","))
	}
	inter := obj.(definition)
	w.Description = inter.GetDescription()
	w.Name = inter.GetName()
	w.EntryPoint = inter.GetPath()
	return nil
}

func (w *Plug) setCrud(name string, obj interface{}) []string {
	unimplimented := checkUnimplimented(obj, (*Crud)(nil))
	if len(unimplimented) > 0 {
		w.Crud = nil
		return unimplimented
	}
	inter := obj.(Crud)
	w.Crud = &inter
	return unimplimented
}

func (w *Plug) setModel(p *plugin.Plugin) []string {

	// Check for Models
	expected := []string{"Create", "Modify", "Delete", "Read"}
	m := make(map[string](func(map[string]string) interface{}), len(expected))
	missing := []string{}
	for _, fun := range expected {
		obj, err := p.Lookup(fun)
		if err != nil {
			missing = append(missing, fun)
			continue
		}
		m[fun] = obj.(func(map[string]string) interface{})
	}
	if len(missing) > 0 {
		return missing
	}
	w.Model = m
	return missing
}
