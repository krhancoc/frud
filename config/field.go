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

func (field *Field) validateType(extraTypes map[string]string) error {
	for key, val := range extraTypes {
		switch f := field.ValueType.(type) {
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
