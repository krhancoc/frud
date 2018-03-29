package config_test

import (
	"testing"

	"github.com/krhancoc/frud/config"
)

var req = &config.DBRequest{
	Method: "post",
	Type:   "meeting",
	Model: config.Fields{
		&config.Field{
			Key:       "another",
			ValueType: "string",
			Options:   []string{"id"},
		},
		&config.Field{
			Key:       "emptyFieldAllowed",
			ValueType: "int",
			Options:   []string{"empty"},
		},
		&config.Field{
			Key: "subFields",
			ValueType: []*config.Field{
				&config.Field{
					Key:       "another",
					ValueType: "string",
				},
				&config.Field{
					Key:       "emptyFieldAllowed",
					ValueType: "int",
					Options:   []string{"empty"},
				},
				&config.Field{
					Key:        "supervisor",
					ValueType:  "person",
					ForeignKey: "name",
					Options:    []string{"empty"},
				},
			},
		},
	},
}

var testCases = []struct {
	params   map[string]interface{}
	expected bool
}{
	{
		map[string]interface{}{
			"another": "anotherthing",
			"subFields": map[string]interface{}{
				"another":           "yep",
				"emptyFieldAllowed": 3,
				"supervisor":        "bob",
			},
		}, true,
	},
	{
		map[string]interface{}{
			"another":   "anotherthing",
			"subFields": map[string]interface{}{},
		}, false,
	},
	{
		map[string]interface{}{
			"another": "anotherthing",
		}, false,
	},
	{
		map[string]interface{}{
			"another": "anotherthing",
			"subFields": map[string]interface{}{
				"another":           "yep",
				"emptyFieldAllowed": "BAD_VALUE",
			},
		}, false,
	},
	{
		map[string]interface{}{
			"another": 420,
			"subFields": map[string]interface{}{
				"another":           "yep",
				"emptyFieldAllowed": "BAD_VALUE",
			},
		}, false,
	},
}

func TestDBRequest(t *testing.T) {
	for i, test := range testCases {
		req.Params = test.params
		err := req.Validate()
		if err != nil && test.expected {
			t.Errorf(err.Error())
		} else if err == nil && !test.expected {
			t.Errorf("Error expected but passed - Test #%d", i)
		}
	}

}
