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

func (c *Cypher) Create() *Command {
	return &Command{
		Type:    "CREATE",
		Req:     c.Req,
		storage: c.Statements,
		Vars:    c.Vars,
	}
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
	return c.paddNext(command)

}

func (c *Cypher) Delete() *Cypher {
	return c.end("DETACH DELETE")
}

func (c *Cypher) Return() *Cypher {
	println("RETURNING")
	return c.end("RETURN")
}

func (c *Cypher) end(t string) *Cypher {

	id := c.Req.Model.GetId()
	values := c.Req.Params
	if t == "RETURN" {
		values = c.Req.Queries
	}
	for _, command := range c.Statements {
		f := command.findVariable(c.Req.Type, id, values[id])
		if f != 0 {
			command := &Command{
				Type:       t,
				Statements: []interface{}{Variable(f)},
				storage:    c.Statements,
				Vars:       c.Vars,
				Req:        c.Req,
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
	return &Cypher{
		Req:        req,
		Statements: []*Command{},
		Vars:       0,
	}
}
