package config

import "fmt"

// Fields is the array version of Field
type Fields []*Field

// ToMap converts the fields object into a map with each fields key as the map key.
func (f Fields) ToMap() map[string]interface{} {
	m := make(map[string]interface{}, len(f))
	for _, field := range f {
		m[field.Key] = field.ValueType
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
func (f Fields) ForeignKeys() Fields {
	keys := []*Field{}
	for _, val := range f {
		if val.ForeignKey != "" {
			keys = append(keys, val)
		}
	}
	return keys
}

// Atomic will grab all fields of atomic types, meaning - string, int, int64, etc.
func (f Fields) Atomic() Fields {
	keys := []*Field{}
	for _, val := range f {
		if val.ForeignKey == "" {
			keys = append(keys, val)
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
