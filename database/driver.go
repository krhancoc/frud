// Package database will hold all the drivers and the CreateDatabase method.
// TODO: Possibly move DBRequest into the database package
package database

import (
	"fmt"

	"github.com/krhancoc/frud/config"
	"github.com/krhancoc/frud/database/neo"
)

// CreateDatabase will instantiate a database driver based on the configuration provided.
func CreateDatabase(conf config.Configuration) (config.Driver, error) {
	switch conf.Database.Type {
	case "neo4j":
		return neo.CreateNeo(conf)
	default:
		return nil, fmt.Errorf("Unrecognized database")
	}
}
