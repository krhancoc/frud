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

type Relation struct {
	Base         byte
	RelationName string
	Head         byte
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

func (c *Cypher) Create() *Command {
	return &Command{
		Type:    "CREATE",
		Req:     c.Req,
		storage: c.Statements,
		Vars:    c.Vars,
	}
}
func (m *Command) findVariable(t string, id string, value string) byte {
	for _, i := range m.Statements {
		stmt, ok := i.(*Statement)
		if !ok {
			continue
		}
		f := stmt.findVariable(t, id, value)
		if f != 0 {
			return f
		}
	}
	return 0
}

func (m *Statement) findVariable(t string, id string, value string) byte {
	if t == m.Label {
		if v, ok := m.Iden[id]; ok && v == value {
			return m.Variable
		}
	}
	return 0
}

func (r *Relation) String() string {
	return fmt.Sprintf(`(%c)-[:%s]->(%c)`, r.Base, r.RelationName, r.Head)
}

func (c *Cypher) Relations() *Cypher {

	var mainVar byte

	for _, command := range c.Statements {
		id := c.Req.Model.GetId()
		found := command.findVariable(c.Req.Type, id, c.Req.Params[id])
		if found != 0 {
			mainVar = found
			break
		}
	}

	relations := []interface{}{}
	for _, val := range c.Req.Model.ForeignKeys() {
		if v, ok := c.Req.Params[val.Key]; ok {
			for _, command := range c.Statements {
				found := command.findVariable(val.ValueType, val.ForeignKey, v)
				if found != 0 {
					relations = append(relations, &Relation{
						Base:         mainVar,
						RelationName: val.Key,
						Head:         found,
					})
					break
				}
			}
		}
	}
	command := &Command{
		Type:       "CREATE",
		Statements: relations,
		storage:    c.Statements,
		Req:        c.Req,
		Vars:       c.Vars,
	}

	if len(c.Statements) > 0 {
		chars := []interface{}{}
		i := 0
		for i = 0; i < c.Vars; i++ {
			chars = append(chars, string(characters[i]))
		}
		with := &Command{
			Type:       "WITH",
			Statements: chars,
		}
		return &Cypher{
			Req:        c.Req,
			Vars:       c.Vars,
			Statements: append(c.Statements, with, command),
		}
	}
	return &Cypher{
		Req:        c.Req,
		Vars:       c.Vars,
		Statements: append(c.Statements, command),
	}

}

func (m *Command) Params() *Cypher {

	stmts := []interface{}{}
	newVars := m.Vars
	for _, val := range m.Req.Model.Atomic() {
		if v, ok := m.Req.Params[val.Key]; ok {
			stmts = append(stmts, &Statement{
				Variable: characters[newVars],
				Label:    m.Req.Type,
				Iden: map[string]interface{}{
					val.Key: v,
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
		stmts = append(stmts, fmt.Sprintf(`%s:"%v"`, key, val))
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

func (c *Cypher) String() string {
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
