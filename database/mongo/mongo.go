package mongo

import (
	"fmt"
	"strings"

	"github.com/krhancoc/frud/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

type Mongo struct {
	Plugins    []*config.PlugConfig
	Connection *mgo.Session
}

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

func (db *Mongo) MakeRequest(req *config.DBRequest) (interface{}, error) {

	session := db.Connection.Clone()
	defer session.Close()
	collection := session.DB("FRUD").C(req.Type)
	switch strings.ToLower(req.Method) {
	case "post":
		req.Params["_id"] = req.Params[req.Model.GetID()]
		err := collection.Insert(req.Params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	case "get":
		var results []map[string]interface{}
		collection.Find(req.Queries).All(&results)
		return results, nil
	default:
		return nil, nil
	}
	return nil, nil
}

func (db *Mongo) ConvertToDriverError(err error) error {
	return nil
}

// ConvertToDriverError(error) error
