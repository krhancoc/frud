package plug

import (
	"fmt"
	"go/importer"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/krhancoc/frud/config"
)

type PlugManager struct {
	Plugs []*Plug
}

type Plug struct {
	Name        string
	Description string
	EntryPoint  string
	Package     *types.Package
	Main        *plugin.Symbol
}

type Definition interface {
	GetName() string
	GetPath() string
	GetDescription() string
}

func getImportString(conf config.PlugConfig) string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir[len(os.Getenv("GOPATH")+"/src/"):] + "/" + conf.PathToCode
}

func plugPack(plug string, conf config.PlugConfig) (*plugin.Plugin, *types.Package, error) {

	importString := getImportString(conf)
	p, err := plugin.Open(conf.PathToCompiled + plug + ".so")
	if err != nil {
		return nil, nil, err
	}
	pkg, err := importer.For("source", nil).Import(importString)
	if err != nil {
		return nil, nil, err
	}
	return p, pkg, nil
}

func isPlugin(o types.Object) bool {

	prefix := strings.HasPrefix(o.Type().Underlying().String(), "struct")
	exported := o.Exported()
	return prefix && exported
}

// CreatePlugManager creates the manager object for Plugins
func CreatePlugManager(conf config.PlugConfig) (*PlugManager, error) {

	var plugManager *PlugManager
	var plugs []*Plug

	for _, plug := range conf.Names {

		p, pkg, err := plugPack(plug, conf)
		if err != nil {
			return plugManager, err
		}

		// Get Package Definition
		for _, name := range pkg.Scope().Names() {
			definition := pkg.Scope().Lookup(name)
			// Check for structure
			if isPlugin(definition) {
				obj, _ := p.Lookup(name)
				switch inter := obj.(type) {
				case Definition:
					thisPlug := Plug{
						Name:        inter.GetName(),
						Description: inter.GetDescription(),
						EntryPoint:  inter.GetPath(),
						Package:     pkg,
						Main:        &obj,
					}
					plugs = append(plugs, &thisPlug)
				default:
					return plugManager, fmt.Errorf("Definition functions not implemented for - %s", name)
				}
			}
		}
	}
	plugManager.Plugs = plugs
	return plugManager, nil
}
