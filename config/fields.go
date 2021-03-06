package config

import (
	"fmt"

	"github.com/krhancoc/frud/errors"
)

// Fields is the array version of Field
type Fields []*Field

// Validate the params with the fields object.  This will traverse the params with the field, to make sure
// that rules are followed
func (f Fields) validateParams(params map[string]interface{}, enforceEmpty bool) error {
	for _, field := range f {
		emptyOpt := field.IsOptionSet("empty")
		v, ok := params[field.Key]
		// If user set this field
		if ok {
			err := field.validate(v, enforceEmpty)
			if err != nil {
				return err
			}
			// Param isnt set, and field is not allowed to be empty.
		} else if !emptyOpt && enforceEmpty {
			return errors.ValidationError{fmt.Sprintf("Required field %s missing", field.Key)}
		}
	}
	return nil
}

func (f Fields) FindField(key string) *Field {
	for _, field := range f {
		if field.Key == key {
			return field
		}
	}
	return nil
}

// ToMap converts the fields object into a map with each fields key as the map key.
func (f Fields) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	for _, field := range f {
		subfields, ok := field.ValueType.(Fields)
		if ok {
			m[field.Key] = subfields.ToMap()
		} else {
			m[field.Key] = field.ValueType
		}
	}
	return m
}

// GetID retrieves the key of the field with the "id" option set.  This option can only be set
// On one field on any given data model.
func (f Fields) GetID() string {
	for _, field := range f {
		for _, option := range field.Options {
			if option == "id" {
				return field.Key
			}
		}
	}
	return ""
}

// ForeignKeys will retrieve all the fields that are not atomic types but rather are types
// that are defined within the model configuration itself.
func (f Fields) ForeignKeys() map[string]*Field {
	keys := make(map[string]*Field)
	for _, val := range f {
		subField, ok := val.ValueType.([]*Field)
		if ok {
			for k, v := range Fields(subField).ForeignKeys() {
				keys[val.Key+"-"+k] = v
			}
		} else {
			if val.ForeignKey != "" {
				keys[val.Key] = val
			}
		}
	}
	return keys
}

// Atomic will grab all fields of atomic types, meaning - string, int, int64, etc.
func (f Fields) Atomic() map[string]*Field {
	keys := make(map[string]*Field)
	for _, val := range f {
		subField, ok := val.ValueType.(Fields)
		if ok {
			for k, v := range subField.Atomic() {
				keys[val.Key+"-"+k] = v
			}
		} else {
			if val.ForeignKey == "" {
				keys[val.Key] = val
			}
		}
	}
	return keys

}

func (f *Fields) validate(extraTypes map[string]string, name string, subfield bool) error {
	idFound := false
	m := make(map[string]bool, len(*f))
	for _, field := range *f {
		if field.Key == "" {
			return errors.ValidationError{
				fmt.Sprintf(`Missing "key" field for a model object in plugin %s`, name),
			}
		}
		if _, ok := m[field.Key]; ok {
			return errors.ValidationError{
				fmt.Sprintf(`Duplicate key - %s - value in model for plugin %s`, field.Key, name),
			}
		}
		err := field.validateType(extraTypes, true)
		if err != nil {
			return err
		}
		m[field.Key] = true
		if field.IsOptionSet("id") {
			if subfield {
				return errors.ValidationError{fmt.Sprintf("Id found in subfield")}
			}
			if idFound {
				return errors.ValidationError{fmt.Sprintf("Multiple id's found in model %s", name)}
			}
			idFound = true
		}
	}
	if !idFound && !subfield {
		return errors.ValidationError{fmt.Sprintf("No id found in %s", name)}
	}
	return nil
}
