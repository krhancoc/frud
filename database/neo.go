package database

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/messages"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/errors"
	"github.com/krhancoc/frud/config"
	log "github.com/sirupsen/logrus"
)

const (
	ConstraintFailure = "Neo.ClientError.Schema.ConstraintValidationFailed"
)

type Neo struct {
	Conf    *config.Database
	Plugins []*config.PlugConfig
}

type constraint struct {
	Model string
	Field string
}

func CreateNeo(conf config.Configuration) (config.Driver, error) {

	var driver config.Driver
	neo := &Neo{
		Conf:    conf.Database,
		Plugins: conf.Manager.Plugs,
	}

	err := neo.initModels()
	if err != nil {
		return driver, err
	}
	var temp interface{} = neo
	driver = temp.(config.Driver)
	return driver, err

}
func (db *Neo) createConstraints(cons []*constraint) error {
	conn := db.Connect().(bolt.Conn)
	defer conn.Close()
	for _, c := range cons {
		stmt := fmt.Sprintf(`CREATE CONSTRAINT ON (n:%s) ASSERT n.%s IS UNIQUE`, c.Model, c.Field)
		log.WithFields(log.Fields{
			"model":    c.Model,
			"field":    c.Field,
			"statment": stmt,
		}).Info("Creating unique constraint")
		n, err := conn.PrepareNeo(stmt)
		if err != nil {
			return err
		}
		_, err = n.ExecNeo(nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Neo) initModels() error {

	log.Infof("Validating models for neo4j database")
	var constraints []*constraint
	idFound := false
	for _, plug := range db.Plugins {
		if plug.Model != nil {
			for _, field := range plug.Model {
				for _, option := range field.Options {
					if option == "id" {
						if idFound {
							return fmt.Errorf("Multiple id's found in model %s", plug.Name)
						}
						idFound = true
						constraints = append(constraints, &constraint{
							Model: plug.Name,
							Field: field.Key,
						})
					}
				}
			}
		}
	}
	return db.createConstraints(constraints)
}

func logHelper(req *config.DBRequest) log.Fields {
	return map[string]interface{}{
		"method": req.Method,
		"type":   req.Type,
		"values": req.Values,
	}
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

	log.WithFields(logHelper(req)).Info("Database request made")
	conn := db.Connect().(bolt.Conn)
	defer conn.Close()

	switch req.Method {
	case "post":
		if Validate(req.Values, req.Model) {
			stmt := fmt.Sprintf(`CREATE (n: %s { %s })`, req.Type, makeValStmt(req.Values, req.Model))
			log.
				WithField("statement", stmt).
				WithFields(logHelper(req)).
				Info("Statement created")
			stmtPrepared, err := conn.PrepareNeo(stmt)
			if err != nil {
				return err
			}
			_, err = stmtPrepared.ExecNeo(nil)
			if err != nil {
				e := err.(*errors.Error).InnerMost().(messages.FailureMessage)
				log.WithFields(e.Metadata).Infof("Attempted post failure")
				if val, ok := e.Metadata["code"]; ok && val == ConstraintFailure {
					return DriverError{
						Status:  http.StatusConflict,
						Message: "Action would violate constraint set by model",
					}
				}
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
