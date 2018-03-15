package config_test

import (
	"testing"

	"github.com/krhancoc/frud/config"
)

var configTests = []struct {
	path  string
	check bool
}{
	{"good.json", true},
	{"no_model_name_2.json", false},
	{"no_model_name.json", false},
	{"replicated_names.json", false},
	{"I_DONT_EXIST.json", false},
	{"unmarshal_error.json", false},
	{"missing_key_field.json", false},
	{"missing_value_type.json", false},
	{"missing_path.json", false},
	{"multiple_ids.json", false},
	{"duplicate_keys.json", false},
}

func TestLoadConfig(t *testing.T) {

	for _, c := range configTests {
		_, err := config.LoadConfig("test_configs/" + c.path)
		if err != nil && c.check {
			t.Errorf("Config failed when should have passed %s", c.path)
		} else if err == nil && !c.check {
			t.Errorf("Config passed when should have failed %s", c.path)
		}
	}
}