package neo

import (
	"fmt"
	"strings"
)

type Relation struct {
	Base         byte
	RelationName string
	Head         byte
}

func (r *Relation) String() string {
	newRelation := strings.Join(strings.Split(r.RelationName, "-"), "_")
	return fmt.Sprintf(`(%c)-[:%s]->(%c)`, r.Base, newRelation, r.Head)
}
