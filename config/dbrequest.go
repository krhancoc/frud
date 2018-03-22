package config

import (
	"fmt"
	"strconv"
)

type DBRequest struct {
	Method  string
	Params  map[string]string
	Queries map[string]string
	Type    string
	Model   Fields
}

func (req *DBRequest) Validate() error {
	for _, field := range req.Model {
		if val, ok := req.Params[field.Key]; ok {
			switch field.ValueType {
			case "int":
				_, err := strconv.ParseInt(val, 10, 32)
				if err != nil {
					return fmt.Errorf("Cannot convert to value type %s", val)
				}
			default:
				continue
			}
		}
	}
	return nil
}

func (req *DBRequest) FollowsModel() error {

	fieldMap := req.Model.ToMap()
	for key := range req.Params {
		if _, ok := fieldMap[key]; !ok {
			return fmt.Errorf("Key %s does not exist within the model for %s", key, req.Type)
		}
	}
	for key := range req.Queries {
		if _, ok := fieldMap[key]; !ok {
			return fmt.Errorf("Key %s does not exist within the model for %s", key, req.Type)
		}
	}
	return nil
}
