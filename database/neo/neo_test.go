package neo_test

import (
	"testing"

	"github.com/krhancoc/frud/config"
	"github.com/krhancoc/frud/database/neo"
)

func TestPostStatement(t *testing.T) {
	req := &config.DBRequest{
		Method: "post",
		Type:   "test",
		Model: []*config.Field{
			&config.Field{
				Key:        "attending",
				ValueType:  "person",
				ForeignKey: "name",
			},
			&config.Field{
				Key:        "meetings",
				ValueType:  "meeting",
				ForeignKey: "date",
			},
			&config.Field{
				Key:       "another",
				ValueType: "string",
			},
		},
		Params: map[string]string{
			"attending": "ken",
			"meetings":  "datehere",
			"another":   "anotherthing",
		},
	}
	println(neo.MakePostStatement(req))
	t.Fail()
}
