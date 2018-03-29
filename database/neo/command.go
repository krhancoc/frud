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

func (m *Command) subStatements(parent string, iden map[string]interface{}, submap map[string]interface{}) {
	depth := strings.Split(parent, "_")
	context := m.Req.Model.ToMap()
	for _, value := range depth {
		context = config.Fields(context[value].([]*config.Field)).ToMap()
	}
	for key, val := range submap {
		if _, ok := context[key]; ok {
			newKey := parent + "_" + key
			subsubmap, ok := val.(map[string]interface{})
			if ok {
				m.subStatements(newKey, iden, subsubmap)
			} else {
				iden[newKey] = val
			}
		}
	}
	return
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
	iden := make(map[string]interface{})

	for key, val := range m.Req.Model.Atomic() {
		if v, ok := requestValues[key]; ok {
			submap, ok := v.(map[string]interface{})
			if ok {
				m.subStatements(val.Key, iden, submap)
			} else {
				iden[key] = v
			}
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
	for key, val := range m.Req.Model.ForeignKeys() {
		if v, ok := m.Req.Params[key]; ok {
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

func (m *Command) findVariable(t string, value map[string]interface{}) byte {
	for _, i := range m.Statements {
		stmt, ok := i.(*Statement)
		if !ok {
			continue
		}
		f := stmt.findVariable(t, value)
		if f != 0 {
			return f
		}
	}
	return 0
}

func setHelper(i interface{}) string {
	stmt := i.(*Statement)
	var equals []string
	for key, val := range stmt.Iden {
		equals = append(equals, fmt.Sprintf(`%c.%s = "%s"`, stmt.Variable, key, val))
	}
	return strings.Join(equals, ",")
}

func (c *Command) String() string {
	stmts := []string{}
	for _, statement := range c.Statements {
		var str string
		if c.Type == "SET" {
			str = fmt.Sprintf("%s", setHelper(statement))
		} else {
			str = fmt.Sprintf("%v", statement)
		}
		stmts = append(stmts, str)
	}

	return c.Type + " " + strings.Join(stmts, ",")
}
