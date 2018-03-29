package neo

import (
	"fmt"
	"strings"

	"github.com/krhancoc/frud/config"
)

//TODO: Still doesnt recognize other types - only strings currently.  Need to bring in creating cypher with proper type.
var characters = "abcdefghijklmnopqrxtuv"

type Variable byte

func (v Variable) String() string {
	return fmt.Sprintf("(%c)", v)
}

type Cypher struct {
	Req        *config.DBRequest
	Statements []*Command
	Vars       int
}

func (c *Cypher) paddNext(command *Command) *Cypher {
	if len(command.Statements) == 0 {
		return &Cypher{
			Req:        c.Req,
			Vars:       c.Vars,
			Statements: c.Statements,
		}
	}
	if len(c.Statements) > 0 {
		return &Cypher{
			Req:        c.Req,
			Vars:       c.Vars,
			Statements: append(c.Statements, CreateWith(c.Vars), command),
		}
	}
	return &Cypher{
		Req:        c.Req,
		Vars:       c.Vars,
		Statements: append(c.Statements, command),
	}
}

func (c *Cypher) Match() *Command {
	return &Command{
		Type:    "MATCH",
		Req:     c.Req,
		storage: c.Statements,
		Vars:    c.Vars,
	}
}

func (c *Cypher) MatchID() *Cypher {
	id := c.Req.Model.GetID()
	if val, ok := c.Req.Params[id]; ok {
		return &Cypher{
			Req: c.Req,
			Statements: append(c.Statements, &Command{
				Type: "MATCH",
				Req:  c.Req,
				Statements: []interface{}{
					&Statement{
						Variable: characters[c.Vars],
						Label:    c.Req.Type,
						Iden: map[string]interface{}{
							id: val,
						},
					},
				},
			}),
			Vars: c.Vars + 1,
		}
	}
	return nil
}

func (c *Cypher) findVariable(t string, values map[string]interface{}) byte {
	for _, command := range c.Statements {
		if f := command.findVariable(t, values); f != 0 {
			return f
		}
	}
	return 0
}

func (c *Cypher) Set() *Cypher {
	id := c.Req.Model.GetID()
	iden := make(map[string]interface{}, len(c.Req.Params))
	for _, fields := range c.Req.Model.Atomic() {
		if val, ok := c.Req.Params[fields.Key]; ok && fields.Key != id {
			iden[fields.Key] = val
		}
	}

	return &Cypher{
		Req: c.Req,
		Statements: append(c.Statements, &Command{
			Type: "SET",
			Req:  c.Req,
			Statements: []interface{}{
				&Statement{
					Variable: c.findVariable(c.Req.Type, c.appValues()),
					Label:    c.Req.Type,
					Iden:     iden,
				},
			},
			Vars: c.Vars,
		}),
	}
}

func flatten(m map[string]interface{}) map[string]interface{} {
	params := make(map[string]interface{})
	for key, val := range m {
		newMap, ok := val.(map[string]interface{})
		if ok {
			for k, v := range flatten(newMap) {
				params[key+"-"+k] = v
			}
		} else {
			params[key] = val
		}
	}
	return params
}

func (c *Cypher) Create() *Command {
	return &Command{
		Type:    "CREATE",
		Req:     c.Req,
		storage: c.Statements,
		Vars:    c.Vars,
	}
}

func (c *Cypher) appValues() map[string]interface{} {
	id := c.Req.Model.GetID()
	appValues := make(map[string]interface{})
	for key, val := range c.Req.Params {
		if strings.HasPrefix(key, id) {
			appValues[key] = val
		}
	}
	for key, val := range c.Req.Queries {
		if strings.HasPrefix(key, id) {
			appValues[key] = val
		}
	}
	return appValues
}
func (c *Cypher) Relations() *Cypher {

	var mainVar byte
	for _, command := range c.Statements {
		found := command.findVariable(c.Req.Type, c.appValues())
		if found != 0 {
			mainVar = found
			break
		}
	}

	relations := []interface{}{}
	for key, val := range c.Req.Model.ForeignKeys() {
		if v, ok := c.Req.Params[key]; ok {
			for _, command := range c.Statements {
				found := command.findVariable(val.ValueType.(string), map[string]interface{}{
					val.ForeignKey: v,
				})
				if found != 0 {
					relations = append(relations, &Relation{
						Base:         mainVar,
						RelationName: key,
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
	return c.paddNext(command)

}

func (c *Cypher) Delete() *Cypher {
	return c.end("DETACH DELETE")
}

func (c *Cypher) Return() *Cypher {
	return c.end("RETURN")
}

func (c *Cypher) end(t string) *Cypher {

	for _, command := range c.Statements {

		if command.Type == "MATCH" {
			f := command.Statements[0].(*Statement).Variable
			command := &Command{
				Type: t,
				Statements: []interface{}{
					Variable(f),
				},
				storage: c.Statements,
				Vars:    c.Vars,
				Req:     c.Req,
			}
			return c.paddNext(command)
		}
	}
	return c
}

func (c *Cypher) String() string {
	commands := []string{}
	for _, command := range c.Statements {
		commands = append(commands, fmt.Sprintf("%v", command))
	}
	return strings.TrimSpace(strings.Join(commands, " "))
}

func CreateCypher(req *config.DBRequest) *Cypher {
	req.Params = flatten(req.Params)
	return &Cypher{
		Req:        req,
		Statements: []*Command{},
		Vars:       0,
	}
}
