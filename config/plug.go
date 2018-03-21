package config

import (
	"fmt"
)

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

	err := conf.Model.validate(extraTypes, conf.Name)
	if err != nil {
		return err
	}

	return nil

}
