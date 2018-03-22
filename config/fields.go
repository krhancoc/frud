package config

import "fmt"

type Fields []*Field

func (f Fields) ToMap() map[string]string {
	m := make(map[string]string, len(f))
	for _, field := range f {
		m[field.Key] = field.ValueType
	}
	return m
}

func (f Fields) GetId() string {
	for _, field := range f {
		for _, option := range field.Options {
			if option == "id" {
				return field.Key
			}
		}
	}
	return ""
}

func (f Fields) ForeignKeys() Fields {
	keys := []*Field{}
	for _, val := range f {
		if val.ForeignKey != "" {
			keys = append(keys, val)
		}
	}
	return keys
}

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
