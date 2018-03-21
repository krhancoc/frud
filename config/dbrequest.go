package config

import "fmt"

type DBRequest struct {
	Method string
	Values map[string]string
	Type   string
	Model  Fields
}

func (req *DBRequest) FollowsModel() error {

	fieldMap := req.Model.ToMap()
	for key := range req.Values {
		if _, ok := fieldMap[key]; !ok {
			return fmt.Errorf("Key %s does not exist within the model for %s", key, req.Type)
		}
	}
	return nil
}
