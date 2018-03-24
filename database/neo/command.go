package neo

import (
	"fmt"
	"strings"

	"github.com/krhancoc/frud/config"
)

type Command struct {
	Type       string
	Req        *config.DBRequest
	Statements []interface{}
	storage    []*Command
	Vars       int
}

func (m *Command) paddNext(newVars int) *Cypher {
	if len(m.Statements) == 0 {
		return &Cypher{
			Req:        m.Req,
			Vars:       newVars,
			Statements: m.storage,
		}
	}
	if len(m.storage) > 0 {
		return &Cypher{
			Req:        m.Req,
			Vars:       newVars,
			Statements: append(m.storage, CreateWith(m.Vars), m),
		}
	}
	return &Cypher{
		Req:        m.Req,
		Vars:       newVars,
		Statements: append(m.storage, m),
	}

}

func (m *Command) Params() *Cypher {
	return m.context("params")
}

func (m *Command) Queries() *Cypher {
	return m.context("queries")
}

func (m *Command) context(values string) *Cypher {

	var requestValues map[string]interface{}
	if values == "params" {
		requestValues = m.Req.Params
	} else {
		requestValues = m.Req.Queries
	}

	stmts := []interface{}{}
	newVars := m.Vars
	iden := make(map[string]interface{}, len(m.Req.Model.Atomic()))

	for _, val := range m.Req.Model.Atomic() {
		if v, ok := requestValues[val.Key]; ok {
			iden[val.Key] = v
		}
	}
	if len(iden) > 0 {
		stmts = append(stmts, &Statement{
			Variable: characters[newVars],
			Label:    m.Req.Type,
			Iden:     iden,
		})
		newVars++
	}
	m.Statements = stmts
	return m.paddNext(newVars)
}

func (m *Command) ForeignKeys() *Cypher {

	stmts := []interface{}{}
	newVars := m.Vars
	for _, val := range m.Req.Model.ForeignKeys() {
		if v, ok := m.Req.Params[val.Key]; ok {
			stmts = append(stmts, &Statement{
				Variable: characters[newVars],
				Label:    val.ValueType.(string),
				Iden: map[string]interface{}{
					val.ForeignKey: v,
				},
			})
			newVars++
		}
	}
	m.Statements = stmts
	return m.paddNext(newVars)
}

func (m *Command) findVariable(t string, id string, value interface{}) byte {
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

func (c *Command) String() string {
	stmts := []string{}
	for _, statement := range c.Statements {
		stmts = append(stmts, fmt.Sprintf("%v", statement))
	}
	return c.Type + " " + strings.Join(stmts, ",")
}
