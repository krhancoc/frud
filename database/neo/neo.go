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
		return CreateCypher(req).Match().Queries().Return().String()
	case "put":
		// return CreateCypher(req).Change().String(), nil
	case "delete":
		return CreateCypher(req).
			Match().Params().Delete().String()
	}
	return ""

}

func (db *Neo) MakeRequest(req *config.DBRequest) (interface{}, error) {

	log.WithFields(logHelper(req)).Info("Database request made")
	conn := *db.Connection

	err := req.Validate()
	if err != nil {
		return nil, err
	}

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
		_, err := stmtPrepared.ExecNeo(nil)
		return nil, db.ConvertToDriverError(err)
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
