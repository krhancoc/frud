// Package config is the package that holds all the configurable structures and functions.
// Its main purpose is to allow users to load a configuration from a json file.
package config

import (
	"encoding/json"
	"io/ioutil"
)

// LoadConfig is the main way of interacting and creating a configuration. This will create and validate a
// configuration from a filename.
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
	return configuration, configuration.Validate()
}
