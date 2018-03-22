package database

import (
	"github.com/krhancoc/frud/config"
	"github.com/krhancoc/frud/database/neo"
)

func CreateDatabase(conf config.Configuration) (config.Driver, error) {
	switch conf.Database.Type {
	case "neo4j":
		return neo.CreateNeo(conf)
	default:
		return nil, nil
	}
}
