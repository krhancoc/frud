package config

import (
	"encoding/json"
	"io/ioutil"
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
	return configuration, configuration.Validate()
}
