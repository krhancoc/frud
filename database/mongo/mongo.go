package mongo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/krhancoc/frud/config"
	"github.com/krhancoc/frud/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mongo Struct
type Mongo struct {
	Plugins    []*config.PlugConfig
	Connection *mgo.Session
}

// CreateMongo creates a database driver
func CreateMongo(conf config.Configuration) (config.Driver, error) {

	var driver config.Driver
	var url string
	db := conf.Database
	if db.User == "" {
		url = fmt.Sprintf("mongodb://%s:%d", db.Hostname, db.Port)
	} else {
		url = fmt.Sprintf("mongodb://%s:%s@%s:%d", db.User, db.Password, db.Hostname, db.Port)
	}
	log.WithField("Connection", url).Debug("Creating connection")

	session, err := mgo.Dial(url)
	if err != nil {
		return driver, err
	}
	return &Mongo{
		Plugins:    conf.Manager.Plugs,
		Connection: session,
	}, nil
}

// MakeRequest implementation of the Mongo Driver.
func (db *Mongo) MakeRequest(req *config.DBRequest) (interface{}, error) {

	// Best practices for the MGO library states that its good to clone the connection and then close that one.
	// We will always have 1 master connection which all requests will clone off of.
	session := db.Connection.Clone()
	defer session.Close()

	collection := session.DB("FRUD").C(req.Type)
	switch strings.ToLower(req.Method) {
	case "post":
		req.Params["_id"] = req.Params[req.Model.GetID()]
		err := collection.Insert(req.Params)
		if err != nil {
			return nil, db.ConvertToDriverError(err)
		}
		return nil, nil
	case "delete":
		req.Params["_id"] = req.Params[req.Model.GetID()]
		err := collection.Remove(req.Params)
		if err != nil {
			return nil, db.ConvertToDriverError(err)
		}
		return nil, nil
	case "get":
		// This allows people to search with params or queries.  This helps when searching for things with
		// nested structure.  May support the ability to dive into nested structures later.
		// An implication of this is that if both are set then it will produce a union of the two searches,
		// May have to perform a map merge
		var qresults []map[string]interface{}
		var presults []map[string]interface{}
		collection.Find(req.Queries).All(&qresults)
		collection.Find(req.Params).All(&presults)
		return append(qresults, presults...), nil
	case "put":
		id := req.Params[req.Model.GetID()].(string)
		err := collection.UpdateId(id, bson.M{"$set": req.Params})
		if err != nil {
			return nil, db.ConvertToDriverError(err)
		}
		return nil, nil
	default:
		return nil, nil
	}
}

func (db *Mongo) ConvertToDriverError(err error) error {
	return errors.DriverError{
		Status:  http.StatusBadRequest,
		Message: err.Error(),
	}
}

// ConvertToDriverError(error) error
