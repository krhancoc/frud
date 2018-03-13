package database

import (
	"strconv"

	"github.com/krhancoc/frud/config"
)

func CreateDatabase(conf *config.Database) config.Driver {
	switch conf.Type {
	case "neo4j":
		var neo interface{} = &Neo{
			Conf: conf,
		}
		return neo.(config.Driver)
	default:
		return nil
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
