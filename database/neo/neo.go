package neo

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/errors"
	"github.com/krhancoc/frud/config"
	frudError "github.com/krhancoc/frud/errors"
	log "github.com/sirupsen/logrus"
)

const (
	ConstraintFailure = "Neo.ClientError.Schema.ConstraintValidationFailed"
)

type Neo struct {
	Plugins    []*config.PlugConfig
	Connection *bolt.Conn
}

type constraint struct {
	Model string
	Field string
}

func CreateNeo(conf config.Configuration) (config.Driver, error) {

	var driver config.Driver
	connection, err := bolt.NewDriver().OpenNeo(fmt.Sprintf("bolt://%s:%s@%s:%d",
		conf.Database.User, conf.Database.Password, conf.Database.Hostname, conf.Database.Port))
	if err != nil {
		panic(err)
	}
	neo := &Neo{
		Plugins:    conf.Manager.Plugs,
		Connection: &connection,
	}
	err = neo.initModels()
	if err != nil {
		return driver, err
	}
	var temp interface{} = neo
	driver = temp.(config.Driver)
	return driver, err

}

func makeValStmt(vals map[string]string, model []*config.Field) string {

	var entries []string
	for _, field := range model {
		if val, ok := vals[field.Key]; ok {
			switch field.ValueType {
			case "int":
				i, _ := strconv.ParseInt(val, 10, 32)
				entries = append(entries, fmt.Sprintf(`%s:%d`, field.Key, i))
			default:
				entries = append(entries, fmt.Sprintf(`%s:"%s"`, field.Key, val))
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
	return frudError.DriverError{
		Status:  http.StatusConflict,
		Message: "Problem with request",
	}
}

func (db *Neo) MakeRequest(req *config.DBRequest) (interface{}, error) {

	log.WithFields(logHelper(req)).Info("Database request made")
	conn := *db.Connection

	err := req.Validate()
	if err != nil {
		return nil, err
	}

	var stmt string

	switch strings.ToLower(req.Method) {
	case "post":
		stmt = MakePostStatement(req)
	case "get":
		stmt = fmt.Sprintf(`MATCH (n: %s { %s }) RETURN (n)`, req.Type, makeValStmt(req.Queries, req.Model))
	case "delete":
		vals := makeValStmt(req.Params, req.Model)
		if vals == "" {
			return nil, frudError.DriverError{
				Status:  http.StatusBadRequest,
				Message: "Bad request",
			}
		}
		stmt = fmt.Sprintf(`MATCH (n: %s { %s }) DETACH DELETE (n)`, req.Type, vals)
	case "put":
		err := req.FollowsModel()
		if err != nil {
			return nil, frudError.DriverError{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			}
		}
		id := req.Model.GetId()
		val, ok := req.Params[id]
		if !ok {
			return nil, frudError.DriverError{
				Status:  http.StatusBadRequest,
				Message: "No ID found in request",
			}
		}
		stmt = fmt.Sprintf(`MATCH (n: %s { %s:"%s" }) SET %s`, req.Type, id, val, makePutStmt(req.Params))
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

// func (db *Neo) Connect() interface{} {
// 	driver := bolt.NewDriver()
// 	connection, err := driver.OpenNeo(fmt.Sprintf("bolt://%s:%s@%s:%d",
// 		db.Conf.User, db.Conf.Password, db.Conf.Hostname, db.Conf.Port))
// 	if err != nil {
// 		log.Fatal(err)
// 		panic(err)
// 	}
// 	return connection
// }
