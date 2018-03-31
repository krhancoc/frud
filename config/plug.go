package config

import (
	"fmt"
)

// PlugConfig encapsulates everything around the plugin object, this will hold data for both
// the model method, or code method of setting up your endpoint.
// TODO: Check if this is the best way
type PlugConfig struct {
	PathToCode     string `json:"pathtocode,omitempty"`
	PathToCompiled string `json:"pathtocompiled,omitempty"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Path           string `json:"path,omitempty"`
	Model          Fields `json:"model,omitempty"`
}

func (conf *PlugConfig) validate(extraTypes map[string]string) error {

	if conf.Name == "" {
		return fmt.Errorf(`Missing "name" field for plugin`)
	}
	if conf.Path == "" {
		return fmt.Errorf(`Missing "path" field for plugin`)
	}
	if len(conf.Model) > 0 {
		err := conf.Model.validate(extraTypes, conf.Name, false)
		if err != nil {
			return err
		}
		return nil
	}
	if conf.PathToCode == "" {
		return fmt.Errorf(`Model field not set and path to code not set`)
	}
	if conf.PathToCompiled == "" {
		return fmt.Errorf(`Model field not set and path to compiled not set`)
	}
	return nil

}
