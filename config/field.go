package config

import (
	"encoding/json"
	"fmt"

	"github.com/krhancoc/frud/errors"
)

// Field is the encapsulation of each of the fields of a object within a database. Its the Field of
// a data model.
type Field struct {
	Key        string      `json:"key"`
	ValueType  interface{} `json:"value_type,omitempty"`
	Options    []string    `json:"options,omitempty"`
	ForeignKey string      `json:"foreignkey,omitempty"`
}

func (field *Field) validate(val interface{}, enforceEmpty bool) error {
	f, ok := field.ValueType.([]*Field)
	if ok {
		switch v := val.(type) {
		case map[string]interface{}:
			return Fields(f).validateParams(v, enforceEmpty)
		default:
			return errors.ValidationError{fmt.Sprintf("Not correct type for field %s", field.Key)}
		}
	}

	switch val.(type) {
	case map[string]interface{}:
		return errors.ValidationError{fmt.Sprintf("Not correct type for field %s", field.Key)}
	case int:
		if field.ValueType.(string) != "int" {
			return errors.ValidationError{fmt.Sprintf("Not correct type for field %s", field.Key)}
		}
	default:
		if field.ValueType.(string) != "string" && field.ForeignKey == "" {
			return errors.ValidationError{fmt.Sprintf("Not correct type for field %s", field.Key)}
		}
	}

	return nil
}

// IsOptionSet checks to see if within that field a specific option flag is
// set.
func (field *Field) IsOptionSet(option string) bool {
	for _, o := range field.Options {
		if o == option {
			return true
		}
	}
	return false
}

// Validate the valueType field
func (field *Field) validateType(extraTypes map[string]string, subfield bool) error {

	//Check extra types, these are values set by the use itself.
	for key, val := range extraTypes {
		// If this is a key field, then lets set the foreignKey field for future use
		// This will help us when we need to know the foreign keys field name
		if f, ok := field.ValueType.(string); ok && f == key {
			field.ForeignKey = val
			return nil
		}
	}
	//Check allowed types
	for _, t := range allowedTypes {
		val, ok := field.ValueType.(string)
		if ok && val == t {
			return nil
		}
	}

	interfaces, ok := field.ValueType.([]interface{})
	if ok && len(interfaces) > 0 {
		var fields Fields
		for _, i := range interfaces {

			b, err := json.Marshal(i)
			if err != nil {
				return err
			}

			var f Field
			err = json.Unmarshal(b, &f)
			if err != nil {
				return errors.ValidationError{"Problem converting subfield"}
			}
			fields = append(fields, &f)
		}
		field.ValueType = fields

		return fields.validate(extraTypes, field.Key, true)
	}
	return errors.ValidationError{
		fmt.Sprintf("Could not find type %s, allowed types are %v or %v", field.ValueType, extraTypes, allowedTypes),
	}
}
