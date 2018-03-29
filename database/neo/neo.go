// Package neo holds all the logic around instantiating a Driver interface around the Neo4J database architecture
package neo

import (
	"fmt"
	"net/http"
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

// Neo is the main neo struct hold the plugins and connection object
type Neo struct {
	Plugins    []*config.PlugConfig
	Connection *bolt.Conn
}

type constraint struct {
	Model string
	Field string
}

// CreateNeo will create a Neo Driver from a configuration object
func CreateNeo(conf config.Configuration) (config.Driver, error) {

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
		return neo, err
	}
	return neo, err

}

// ConvertToDriverError will convert a Neo4J connector error to a generic Driver error.
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

func createStatement(req *config.DBRequest) string {

	switch strings.ToLower(req.Method) {
	case "post":
		return CreateCypher(req).
			Match().ForeignKeys().
			Create().Params().
			Relations().
			String()
	case "get":
		return CreateCypher(req).Match().Queries().String() + "-[b]-(c) RETURN a,b,c"
	case "put":
		return CreateCypher(req).MatchID().Set().String()
	case "delete":
		return CreateCypher(req).
			Match().Params().Delete().String()
	}
	return ""

}

func (db *Neo) MakeRequest(req *config.DBRequest) (interface{}, error) {

	log.WithFields(logHelper(req)).Info("Database request made")
	conn := *db.Connection

	stmt := createStatement(req)
	println(stmt)

	log.
		WithField("statement", stmt).
		WithFields(logHelper(req)).
		Info("Statement created")

	stmtPrepared, err := conn.PrepareNeo(stmt)
	defer stmtPrepared.Close()

	if err != nil {
		return nil, db.ConvertToDriverError(err)
	}

	switch strings.ToLower(req.Method) {
	case "post", "delete", "put":
		result, err := stmtPrepared.ExecNeo(nil)
		if err != nil {

			return nil, db.ConvertToDriverError(err)
		}
		num, _ := result.RowsAffected()
		return num, nil
	case "get":
		result, err := stmtPrepared.QueryNeo(nil)
		if err != nil {
			return nil, db.ConvertToDriverError(err)
		}
		r, _, _ := result.All()
		return r, nil

	}

	return nil, frudError.DriverError{
		Status:  http.StatusBadRequest,
		Message: "Bad Request",
	}
}
