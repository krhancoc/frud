package config

import "fmt"

type Field struct {
	Key        string   `json:"key"`
	ValueType  string   `json:"value_type"`
	Options    []string `json:"options,omitempty"`
	ForeignKey string   `json:"foreignkey,omitempty"`
}

var allowedTypes = []string{
	"int",
	"string",
}

func (field *Field) validateType(extraTypes map[string]string) error {
	for key, val := range extraTypes {
		if field.ValueType == key {
			field.ForeignKey = val
			return nil
		}
	}
	for _, t := range allowedTypes {
		if t == field.ValueType {
			return nil
		}
	}
	return fmt.Errorf("Could not find type %s, allowed types are %v or %v", field.ValueType, extraTypes, allowedTypes)
}
