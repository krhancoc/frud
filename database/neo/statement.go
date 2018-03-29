package neo

import (
	"fmt"
	"strings"
)

type Identifiers map[string]interface{}

type Statement struct {
	Variable byte
	Label    string
	Iden     Identifiers
}

func (s Identifiers) String() string {
	stmts := []string{}
	for key, val := range s {
		stmts = append(stmts, fmt.Sprintf(`%s:'%v'`, key, val))
	}
	return strings.Join(stmts, ",")
}

func (m *Statement) findVariable(t string, value map[string]interface{}) byte {
	if t == m.Label {
		for key, _ := range value {
			if _, ok := m.Iden[key]; ok {
				return m.Variable
			}
		}
	}
	return 0
}

func (s *Statement) String() string {
	return fmt.Sprintf(`(%c:%s { %s })`, s.Variable, s.Label, s.Iden.String())
}
