package neo_test

import (
	"testing"

	"github.com/krhancoc/frud/config"
	"github.com/krhancoc/frud/database/neo"
)

var req = &config.DBRequest{
	Method: "post",
	Type:   "meeting",
	Model: []*config.Field{
		&config.Field{
			Key:        "attending",
			ValueType:  "person",
			ForeignKey: "name",
		},
		&config.Field{
			Key:        "nextup",
			ValueType:  "meeting",
			ForeignKey: "date",
		},
		&config.Field{
			Key:       "another",
			ValueType: "string",
			Options:   []string{"id"},
		},
	},
	Params: map[string]string{
		"attending": "ken",
		"nextup":    "datehere",
		"another":   "anotherthing",
	},
	Queries: map[string]string{
		"attending": "ken",
		"nextup":    "datehere",
		"another":   "anotherthing",
	},
}

func TestCypher(t *testing.T) {
	cypher := neo.CreateCypher(req)
	println(cypher.Match().ForeignKeys().Match().ForeignKeys().String())
	println(cypher.Match().ForeignKeys().Create().Params().Relations().String())
	t.Fail()

}
func TestPostStatement(t *testing.T) {
	println(neo.MakePostStatement(req))
	t.Fail()
}
