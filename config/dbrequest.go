package config

import (
	"fmt"
)

type DBRequest struct {
	Method  string
	Params  map[string]interface{}
	Queries map[string]interface{}
	Type    string
	Model   Fields
}

func (req *DBRequest) Validate() error {
	for _, field := range req.Model {
		if val, ok := req.Params[field.Key]; ok {
			switch field.ValueType {
			case "int":
				_, ok := (val).(int)
				if !ok {
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
