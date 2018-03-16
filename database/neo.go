package database

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

	var constraints []*constraint
	for _, plug := range db.Plugins {
		for _, field := range plug.Model {
			for _, option := range field.Options {
				switch option {
				case "id":
					constraints = append(constraints, &constraint{
						Model: plug.Name,
						Field: field.Key,
					})
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
				entries = append(entries, fmt.Sprintf(`%s: %d`, field.Key, i))
			default:
				entries = append(entries, fmt.Sprintf(`%s: "%s"`, field.Key, val))
			}
		}
	}
	return strings.Join(entries, ",")
}

func makePutStmt(vals map[string]string) string {
	stmt := []string{}
	for key, val := range vals {
		stmt = append(stmt, fmt.Sprintf(`n.%s = "%s"`, key, val))
	}
	return strings.Join(stmt, ",")
}

func (db Neo) ConvertToDriverError(err error) error {

	if err == nil {
		return err
	}

	e := err.(*errors.Error).InnerMost()
	log.Error(e.Error())
	return DriverError{
		Status:  http.StatusConflict,
		Message: "Problem with request",
	}
}

func (db *Neo) MakeRequest(req *config.DBRequest) (interface{}, error) {

	log.WithFields(logHelper(req)).Info("Database request made")
	conn := db.Connect().(bolt.Conn)
	defer conn.Close()

	if !Validate(req.Values, req.Model) {
		return "", fmt.Errorf("Invalid request")
	}

	var stmt string

	switch req.Method {
	case "post":
		stmt = fmt.Sprintf(`CREATE (n: %s { %s })`, req.Type, makeValStmt(req.Values, req.Model))
	case "get":
		stmt = fmt.Sprintf(`MATCH (n: %s { %s }) RETURN (n)`, req.Type, makeValStmt(req.Values, req.Model))
	case "delete":
		stmt = fmt.Sprintf(`MATCH (n: %s { %s }) DELETE (n)`, req.Type, makeValStmt(req.Values, req.Model))
	case "put":
		err := req.FollowsModel()
		if err != nil {
			return nil, DriverError{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			}
		}
		id := req.Model.GetId()
		val, ok := req.Values[id]
		if !ok {
			return nil, DriverError{
				Status:  http.StatusBadRequest,
				Message: "No ID found in request",
			}
		}
		stmt = fmt.Sprintf(`MATCH (n: %s { %s:"%s" }) SET %s`, req.Type, id, val, makePutStmt(req.Values))
	}
	log.
		WithField("statement", stmt).
		WithFields(logHelper(req)).
		Info("Statement created")

	stmtPrepared, err := conn.PrepareNeo(stmt)
	defer stmtPrepared.Close()

	if err != nil {
		return nil, db.ConvertToDriverError(err)
	}
	switch req.Method {
	case "post", "delete", "put":
		_, err := stmtPrepared.ExecNeo(nil)
		return nil, db.ConvertToDriverError(err)
	case "get":
		result, err := stmtPrepared.QueryNeo(nil)
		r, _, _ := result.All()
		return r, db.ConvertToDriverError(err)
	}
	return nil, nil
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
