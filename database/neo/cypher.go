package neo

import (
	"fmt"
	"strings"

	"github.com/krhancoc/frud/config"
)

type Identifiers map[string]interface{}

type Command struct {
	Type       string
	Req        *config.DBRequest
	Statements []interface{}
	storage    []*Command
	Vars       int
}

type Statement struct {
	Variable byte
	Label    string
	Iden     Identifiers
}

type Cypher struct {
	Req        *config.DBRequest
	Statements []*Command
	Vars       int
}

var characters = "abcdefghijklmnopqrxtuv"

func (c *Cypher) Match() *Command {
	return &Command{
		Type:    "MATCH",
		Req:     c.Req,
		storage: c.Statements,
		Vars:    c.Vars,
	}
}

func (m *Command) ForeignKeys() *Cypher {
	stmts := []interface{}{}
	newVars := m.Vars
	for _, val := range m.Req.Model.ForeignKeys() {
		if v, ok := m.Req.Params[val.Key]; ok {
			stmts = append(stmts, &Statement{
				Variable: characters[newVars],
				Label:    val.ValueType,
				Iden: map[string]interface{}{
					val.ForeignKey: v,
				},
			})
			newVars++
		}
	}
	m.Statements = stmts
	if len(m.storage) > 0 {
		chars := []interface{}{}
		i := 0
		for i = 0; i < m.Vars; i++ {
			chars = append(chars, string(characters[i]))
		}
		with := &Command{
			Type:       "WITH",
			Statements: chars,
		}
		return &Cypher{
			Req:        m.Req,
			Vars:       newVars,
			Statements: append(m.storage, with, m),
		}
	}
	return &Cypher{
		Req:        m.Req,
		Vars:       newVars,
		Statements: append(m.storage, m),
	}
}

func (s Identifiers) String() string {
	stmts := []string{}
	for key, val := range s {
		stmts = append(stmts, fmt.Sprintf(`%s:%v`, key, val))
	}
	return strings.Join(stmts, ",")
}

func (s *Statement) String() string {
	return fmt.Sprintf(`(%c:%s { %s })`, s.Variable, s.Label, s.Iden.String())
}

func (c *Command) String() string {
	stmts := []string{}
	for _, statement := range c.Statements {
		stmts = append(stmts, fmt.Sprintf("%v", statement))
	}
	return c.Type + " " + strings.Join(stmts, ",")
}

func (c *Cypher) ToString() string {
	commands := []string{}
	for _, command := range c.Statements {
		commands = append(commands, fmt.Sprintf("%v", command))
	}
	return strings.Join(commands, " ")
}

func CreateCypher(req *config.DBRequest) *Cypher {
	return &Cypher{
		Req:        req,
		Statements: []*Command{},
		Vars:       0,
	}
}
