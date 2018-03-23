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
		stmts = append(stmts, fmt.Sprintf(`%s:"%v"`, key, val))
	}
	return strings.Join(stmts, ",")
}

func (m *Statement) findVariable(t string, id string, value string) byte {
	if t == m.Label {
		if v, ok := m.Iden[id]; ok && v == value {
			return m.Variable
		}
	}
	return 0
}

func (s *Statement) String() string {
	return fmt.Sprintf(`(%c:%s { %s })`, s.Variable, s.Label, s.Iden.String())
}
