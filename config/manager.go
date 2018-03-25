package config

import "fmt"

// ManagerConfig is config object that will hold everything todo with the endpoints itself, including
// the datamodels, as well as the generic endpoints that will be set.  Things like healthcheck and eventually,
// the login endpoint
type ManagerConfig struct {
	Generics []string      `json:"generics"`
	Plugs    []*PlugConfig `json:"plugins"`
}

func (conf *ManagerConfig) validateGenerics() error {
	for _, gen := range conf.Generics {
		found := false
		for _, generic := range allowedGenerics {
			if gen == generic {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("No such generic exists - %s", gen)
		}
	}
	return nil
}

func (conf *ManagerConfig) collectTypeIdMap() (map[string]string, error) {
	types := make(map[string]string, len(conf.Plugs))
	for _, plug := range conf.Plugs {
		if _, ok := types[plug.Name]; ok {
			return nil, fmt.Errorf("Duplicate plugin name %s", plug.Name)
		}
		types[plug.Name] = plug.Model.GetID()
	}
	return types, nil
}

func (conf *ManagerConfig) Validate() error {
	err := conf.validateGenerics()
	if err != nil {
		return err
	}
	m, err := conf.collectTypeIdMap()
	if err != nil {
		return err
	}

	paths := make(map[string]bool, len(conf.Plugs))
	for _, plug := range conf.Plugs {
		if _, ok := paths[plug.Path]; ok {
			return fmt.Errorf(`Duplicate path - %s`, plug.Path)
		}
		err = plug.validate(m)
		if err != nil {
			return err
		}
		paths[plug.Path] = true
	}
	return nil
}
