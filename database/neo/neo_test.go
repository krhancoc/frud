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
		&config.Field{
			Key:       "here",
			ValueType: "string",
		},
	},
	Params: map[string]interface{}{
		"attending": "ken",
		"nextup":    "datehere",
		"another":   "anotherthing",
		"here":      "HEY HELLO",
	},
	Queries: map[string]interface{}{
		"attending": "ken",
		"nextup":    "datehere",
		"another":   "anotherthing",
	},
}

var reqTwo = &config.DBRequest{
	Method: "post",
	Type:   "meeting",
	Model: []*config.Field{
		&config.Field{
			Key:       "another",
			ValueType: "string",
			Options:   []string{"id"},
		},
	},
	Params: map[string]interface{}{
		"another": "anotherthing",
	},
}

var reqThree = &config.DBRequest{
	Method: "post",
	Type:   "meeting",
	Model: []*config.Field{
		&config.Field{
			Key: "another",
			ValueType: []*config.Field{
				&config.Field{
					Key:       "hello",
					ValueType: "string",
				},
				&config.Field{
					Key:       "world",
					ValueType: "string",
				},
			},
			Options: []string{"id"},
		},
	},
	Params: map[string]interface{}{
		"another": "anotherthing",
	},
}

func TestCypher(t *testing.T) {
	cypher := neo.CreateCypher(req)
	cypherTwo := neo.CreateCypher(reqTwo)
	cypherThree := neo.CreateCypher(reqThree)
	println(cypher.Match().ForeignKeys().Match().ForeignKeys().String())
	println(cypher.Match().ForeignKeys().Create().Params().Relations().String())
	println(cypherTwo.Create().Params().String())
	println(cypherTwo.Match().ForeignKeys().Create().Params().Relations().String())
	println(cypherTwo.Match().Params().Delete().String())
	println(cypherTwo.Match().Queries().Return().String())
	println(cypherThree.Match().Params().Return().String())
	println(cypher.Match().Params().String())
	println(cypher.MatchID().Set().String())
	t.Fail()

}
func TestPostStatement(t *testing.T) {
	t.Fail()
}
