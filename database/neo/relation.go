package neo

import "fmt"

type Relation struct {
	Base         byte
	RelationName string
	Head         byte
}

func (r *Relation) String() string {
	return fmt.Sprintf(`(%c)-[:%s]->(%c)`, r.Base, r.RelationName, r.Head)
}
