package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func LoadConfig(filename string) Configuration {
	configuration := Configuration{}
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Problem loading configuration file -- make sure its in the root of the project")
		os.Exit(1)
	}
	err = json.Unmarshal(raw, &configuration)
	if err != nil {
		log.Fatal("Error:", err)
	}
	return configuration
}
