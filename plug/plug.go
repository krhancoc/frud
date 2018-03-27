// Package plug encapsulates the plugin style that this framework uses.  Using either the code method,
// which relies on the go plugin package, as well as the package package (yes.) Or the configuration style
// of defining models within a json file.  Plugins are added to a pluginManager which is added to the server.
package plug

import (
	"fmt"
	"go/importer"
	"go/types"
	"plugin"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/krhancoc/frud/config"
)

type endpoint struct {
	Name        string
	Path        string
	Description string
}

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

func getImportString(conf *config.PlugConfig) string {

	return "github.com/krhancoc/frud/" + conf.PathToCode
}

func isPlugin(o types.Object) bool {

	prefix := strings.HasPrefix(o.Type().Underlying().String(), "struct")
	exported := o.Exported()
	return prefix && exported
}

func plugPack(conf *config.PlugConfig) (*plugin.Plugin, *types.Package, error) {

	importString := getImportString(conf)
	p, err := plugin.Open(conf.PathToCompiled)
	if err != nil {
		return nil, nil, err
	}
	pkg, err := importer.For("source", nil).Import(importString)
	if err != nil {
		return nil, nil, err
	}
	return p, pkg, nil
}

func (w *Plug) setDefinition(obj interface{}) {
	data := reflect.ValueOf(obj).Elem().FieldByName("Data").Interface().(*endpoint)
	w.Name = data.Name
	w.Description = data.Description
	w.Path = data.Path
}

func (w *Plug) setCrud(obj interface{}) []string {
	unimplimented := checkUnimplimented(obj, (*Crud)(nil))
	if len(unimplimented) > 0 {
		w.Crud = nil
		return unimplimented
	}
	inter := obj.(Crud)
	w.Crud = &inter
	return unimplimented
}

func createPlug(conf *config.PlugConfig) (*Plug, error) {
	if conf.Model == nil {
		return createPlugFromCode(conf)
	}
	return createPlugFromModel(conf)
}

func createPlugFromModel(conf *config.PlugConfig) (*Plug, error) {
	color.Yellow("Plugin found using Model Method - %s", conf.Name)
	thisPlug := &Plug{
		Name:        conf.Name,
		Description: conf.Description,
		Path:        conf.Path,
		Model:       conf.Model,
		Crud:        nil,
	}
	return thisPlug, nil
}

func createPlugFromCode(conf *config.PlugConfig) (*Plug, error) {
	color.Yellow("Plugin found using Code Method - %s", conf.Name)
	thisPlug := &Plug{
		Name:        "",
		Description: "",
		Path:        "",
		Model:       nil,
		Crud:        nil,
	}
	var unimplemented []string

	p, pkg, err := plugPack(conf)
	if err != nil {
		return thisPlug, err
	}
	// Check for Handlers
	for _, name := range pkg.Scope().Names() {
		definition := pkg.Scope().Lookup(name)
		// Check for structure
		if isPlugin(definition) {
			obj, err := p.Lookup(name)
			if err != nil {
				continue
			}
			thisPlug.setDefinition(obj)
			unimplemented = thisPlug.setCrud(obj)
			if len(unimplemented) > 0 {
				return thisPlug, fmt.Errorf("Unimplimented methods in %s: %s", thisPlug.Name, strings.Join(unimplemented, ","))
			}
		}
	}
	return thisPlug, nil
}
