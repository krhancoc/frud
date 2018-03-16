package config

import "fmt"

func (req *DBRequest) FollowsModel() error {

	fieldMap := req.Model.ToMap()
	for key, _ := range req.Values {
		if _, ok := fieldMap[key]; !ok {
			return fmt.Errorf("Key %s does not exist within the model for %s", key, req.Type)
		}
	}
	return nil
}
