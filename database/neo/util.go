package neo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/krhancoc/frud/config"
	log "github.com/sirupsen/logrus"
)

func (db *Neo) createConstraints(cons []*constraint) error {
	conn := *db.Connection
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
		n.Close()
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
		"method":  req.Method,
		"type":    req.Type,
		"params":  req.Params,
		"queries": req.Queries,
	}
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
