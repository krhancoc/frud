package config

import (
	"fmt"
)

// Fields is the array version of Field
type Fields []*Field

func (f Fields) ValidateParams(params map[string]interface{}) error {
	for _, field := range f {
		empty := field.IsOptionSet("empty")
		v, ok := params[field.Key]
		if ok {
			err := field.Validate(v)
			if err != nil {
				return err
			}
		} else if !empty {
			return fmt.Errorf("Required field %s missing", field.Key)
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

func (f *Fields) validate(extraTypes map[string]string, name string) error {
	idFound := false

	m := make(map[string]bool, len(*f))
	for _, field := range *f {
		if field.Key == "" {
			return fmt.Errorf(`Missing "key" field for a model object in plugin %s`, name)
		}
		if _, ok := m[field.Key]; ok {
			return fmt.Errorf(`Duplicate key - %s - value in model for plugin %s`, field.Key, name)
		}
		err := field.validateType(extraTypes)
		if err != nil {
			return err
		}
		m[field.Key] = true
		for _, option := range field.Options {
			if option == "id" {
				if idFound {
					return fmt.Errorf("Multiple id's found in model %s", name)
				}
				idFound = true
			}
		}
	}
	return nil
}
