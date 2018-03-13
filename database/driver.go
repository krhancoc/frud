package database

import (
	"fmt"
	"log"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/krhancoc/frud/config"
)

func CreateDatabase(conf *config.Database) config.Driver {
	switch conf.Type {
	case "neo4j":
		var neo interface{} = &Neo{
			Conf: conf,
		}
		return neo.(config.Driver)
	default:
		return nil
	}
}

type Neo struct {
	Conf *config.Database
}

func (db *Neo) MakeRequest(method string, vals map[string]string, model []*config.Field) error {

	println("Connecting to", db.Conf.Type, "DB")
	conn := db.Connect().(bolt.Conn)
	conn.Close()
	println("Closing Connection to DB")

	return nil
}

func (db *Neo) Connect() interface{} {
	driver := bolt.NewDriver()
	connection, err := driver.OpenNeo(fmt.Sprintf("bolt://%s:%s@%s:%d",
		db.Conf.User, db.Conf.Password, db.Conf.Hostname, db.Conf.Port))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return connection
}
