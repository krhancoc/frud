package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

func LoadConfig(filename string) (Configuration, error) {
	configuration := Configuration{}
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return configuration, err
	}
	err = json.Unmarshal(raw, &configuration)
	if err != nil {
		return configuration, err
	}
	return configuration, ValidateConfiguration(configuration)
}

func validateModelMethod(conf *PlugConfig) error {
	log.Infof("Validating Plugin")
	idFound := false
	if conf.Name == "" {
		return fmt.Errorf(`Missing "name" field for plugin`)
	}
	if conf.Path == "" {
		return fmt.Errorf(`Missing "path" field for plugin`)
	}
	m := make(map[string]bool, len(conf.Model))
	for _, field := range conf.Model {
		if field.Key == "" {
			return fmt.Errorf(`Missing "key" field for a model object in plugin %s`, conf.Name)
		}
		if field.ValueType == "" {
			return fmt.Errorf(`Missing "value_type" field for a model object in plugin %s`, conf.Name)
		}
		if _, ok := m[field.Key]; ok {
			return fmt.Errorf(`Duplicate key - %s - value in model for plugin %s`, field.Key, conf.Name)
		}
		m[field.Key] = true
		for _, option := range field.Options {
			if option == "id" {
				if idFound {
					return fmt.Errorf("Multiple id's found in model %s", conf.Name)
				}
				idFound = true
			}
		}
	}
	return nil
}

func validatePlugins(conf []*PlugConfig) error {

	names := make(map[string]bool, len(conf))
	paths := make(map[string]bool, len(conf))
	for _, plug := range conf {

		if _, ok := names[plug.Name]; ok {
			return fmt.Errorf(`Duplicate name - %s`, plug.Name)
		}
		names[plug.Name] = true

		if _, ok := names[plug.Path]; ok {
			return fmt.Errorf(`Duplicate path - %s`, plug.Path)
		}
		paths[plug.Path] = true

		if len(plug.Model) > 0 {
			if err := validateModelMethod(plug); err != nil {
				return err
			}
		}
	}
	return nil
}

func ValidateConfiguration(conf Configuration) error {
	return validatePlugins(conf.Manager.Plugs)
}
