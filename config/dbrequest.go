package config

import (
	"fmt"

	"github.com/krhancoc/frud/errors"
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
	switch req.Method {
	case "delete":
		_, ok := req.Params[req.Model.GetID()]
		if !ok {
			return errors.ValidationError{fmt.Sprintf("ID field of model required")}
		}
	case "post":
		// Enforce empty rules
		return req.Model.validateParams(req.Params, true)
	case "put":
		_, ok := req.Params[req.Model.GetID()]
		if !ok {
			return errors.ValidationError{fmt.Sprintf("ID field of model required")}
		}
		// Don't enforce empty rules on put's
		return req.Model.validateParams(req.Params, false)
	case "get":
		return req.Model.validateParams(req.Params, false)
	}
	return nil
}
