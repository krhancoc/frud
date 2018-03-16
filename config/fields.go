package config

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
