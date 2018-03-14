package database

import (
	"strconv"

	"github.com/krhancoc/frud/config"
)

func CreateDatabase(conf config.Configuration) (config.Driver, error) {
	switch conf.Database.Type {
	case "neo4j":
		return CreateNeo(conf)
	default:
		return nil, nil
	}
}

func Validate(vals map[string]string, fields []*config.Field) bool {
	for _, field := range fields {
		if val, ok := vals[field.Key]; ok {
			switch field.ValueType {
			case "int":
				_, err := strconv.ParseInt(val, 10, 32)
				if err != nil {
					return false
				}
			default:
				continue
			}
		}
	}
	return true
}
