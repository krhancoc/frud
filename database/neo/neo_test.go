package neo_test

import (
	"testing"

	"github.com/krhancoc/frud/config"
	"github.com/krhancoc/frud/database/neo"
)

var req = &config.DBRequest{
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
	Queries: map[string]string{
		"attending": "ken",
		"meetings":  "datehere",
		"another":   "anotherthing",
	},
}

func TestCypher(t *testing.T) {
	cypher := neo.CreateCypher(req)
	println(cypher.Match().ForeignKeys().Match().ForeignKeys().ToString())
	t.Fail()

}
func TestPostStatement(t *testing.T) {
	println(neo.MakePostStatement(req))
	t.Fail()
}
