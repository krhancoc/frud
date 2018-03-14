package database

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/krhancoc/frud/config"
)

type Neo struct {
	Conf *config.Database
}

func makeValStmt(vals map[string]string, model []*config.Field) string {

	var entries []string
	for _, field := range model {
		if val, ok := vals[field.Key]; ok {
			switch field.ValueType {
			case "int":
				i, _ := strconv.ParseInt(val, 10, 32)
				entries = append(entries, fmt.Sprintf(`%s: %i`, field.Key, i))
			default:
				entries = append(entries, fmt.Sprintf(`%s: "%s"`, field.Key, val))
			}
		}
	}
	return strings.Join(entries, ",")
}

func (db *Neo) MakeRequest(req *config.DBRequest) error {

	println("Connecting to", db.Conf.Type, "DB")
	conn := db.Connect().(bolt.Conn)
	defer conn.Close()

	switch req.Method {
	case "post":
		if Validate(req.Values, req.Model) {
			stmt := fmt.Sprintf(`CREATE (n: %s { %s })`, req.Type, makeValStmt(req.Values, req.Model))
			println(stmt)
			stmtPrepared, err := conn.PrepareNeo(stmt)
			if err != nil {
				return err
			}
			_, err = stmtPrepared.ExecNeo(nil)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("Invalid request")
	}

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
