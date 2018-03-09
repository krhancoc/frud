package database

import (
	"fmt"
	"log"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/krhancoc/frud/config"
)

type Driver interface {
	Connect(db *config.Database) bolt.Conn
}

func Connect(db *config.Database) bolt.Conn {
	switch db.Type {
	case "neo4j":
		driver := bolt.NewDriver()
		connection, err := driver.OpenNeo(fmt.Sprintf("bolt://%s:%s@%s:%d",
			db.User, db.Password, db.Hostname, db.Port))
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		return connection
	case "mongo":
		return nil
	}
	return nil
}
