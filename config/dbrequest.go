package config

import (
	"fmt"
)

// DBRequest is the object given to our drivers,  each driver must implement the MakeRequest function.
// This function will take a DBRequest object as a parameter. This object will hold a summerized
// version of the Restful call
type DBRequest struct {
	Method  string
	Params  map[string]interface{}
	Queries map[string]interface{}
	Type    string
	Model   Fields
}

// Validate will validate the DBRequest with the model provided to make sure the values are able
// to convert into the proper values given by the model
// TODO: Check/Enforce type conversion on fields - int, int64 etc.
func (req *DBRequest) Validate() error {
	return req.Model.ValidateParams(req.Params)
}

// FollowsModel will check to make sure the DBRequest params and queries follow the model attached to
// the endpoint itself
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
