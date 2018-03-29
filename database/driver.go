// Package database will hold all the drivers and the CreateDatabase method.
// TODO: Possibly move DBRequest into the database package
package database

import (
	"fmt"

	"github.com/krhancoc/frud/config"
	"github.com/krhancoc/frud/database/mongo"
	"github.com/krhancoc/frud/database/neo"
	log "github.com/sirupsen/logrus"
)

// CreateDatabase will instantiate a database driver based on the configuration provided.
func CreateDatabase(conf config.Configuration) (config.Driver, error) {
	log.Debug("Creating database of type " + conf.Database.Type)
	switch conf.Database.Type {
	case "neo4j":
		return neo.CreateNeo(conf)
	case "mongo":
		return mongo.CreateMongo(conf)
	default:
		return nil, fmt.Errorf("Unrecognized database")
	}
}
