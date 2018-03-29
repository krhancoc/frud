package config

import (
	"encoding/json"
	"fmt"
)

// Field is the encapsulation of each of the fields of a object within a database. Its the Field of
// a data model.
type Field struct {
	Key        string      `json:"key"`
	ValueType  interface{} `json:"value_type"`
	Options    []string    `json:"options,omitempty"`
	ForeignKey string      `json:"foreignkey,omitempty"`
}

func (field *Field) Validate(val interface{}) error {
	f, ok := field.ValueType.([]*Field)
	if ok {
		switch v := val.(type) {
		case map[string]interface{}:
			return Fields(f).ValidateParams(v)
		default:
			return fmt.Errorf("Not correct type for field %s", field.Key)
		}
	}

	switch val.(type) {
	case map[string]interface{}:
		return fmt.Errorf("Not correct type for field %s", field.Key)
	case int:
		if field.ValueType.(string) != "int" {
			return fmt.Errorf("Not correct type for field %s", field.Key)
		}
	default:
		if field.ValueType.(string) != "string" && field.ForeignKey == "" {
			return fmt.Errorf("Not correct type for field %s", field.Key)
		}
	}

	return nil
}

func (field *Field) IsOptionSet(option string) bool {
	for _, o := range field.Options {
		if o == option {
			return true
		}
	}
	return false
}

func (field *Field) validateType(extraTypes map[string]string) error {
	for key, val := range extraTypes {
		switch f := field.ValueType.(type) {
		case Fields:
			return f.validate(extraTypes, field.Key)
		case string:
			if f == key {
				field.ForeignKey = val
				return nil
			}
		}
	}
	for _, t := range allowedTypes {
		val, ok := field.ValueType.(string)
		if ok && val == t {
			return nil
		}
	}

	interfaces, ok := field.ValueType.([]interface{})
	if ok && len(interfaces) > 0 {
		var fields []*Field
		for _, i := range interfaces {

			b, err := json.Marshal(i)
			if err != nil {
				return err
			}

			var f Field
			err = json.Unmarshal(b, &f)
			if err != nil {
				return fmt.Errorf("Problem converting subfield")
			}
			err = f.validateType(extraTypes)
			if err != nil {
				return err
			}
			fields = append(fields, &f)
		}
		field.ValueType = fields
		return nil
	}
	return fmt.Errorf("Could not find type %s, allowed types are %v or %v", field.ValueType, extraTypes, allowedTypes)
}
