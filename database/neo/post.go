package neo

import (
	"fmt"
	"strings"

	"github.com/krhancoc/frud/config"
)

func MakePostStatement(req *config.DBRequest) string {

	characters := "abcdefghijklmopqrxtuvwxyz"
	matchStmt := []string{}
	createStmt := []string{}
	charUses := 0
	for _, val := range req.Model.ForeignKeys() {
		if v, ok := req.Values[val.Key]; ok {
			matchStmt = append(matchStmt, fmt.Sprintf(`(%c:%s {%s:"%s"})`,
				characters[charUses], val.ValueType, val.ForeignKey, v))
			createStmt = append(createStmt, fmt.Sprintf(`(n)-[:%s]->(%c)`, val.Key, characters[charUses]))
			charUses++
		}
	}
	vals := []string{}
	for _, val := range req.Model.Atomic() {
		vals = append(vals, fmt.Sprintf(`%s:"%s"`, val.Key, req.Values[val.Key]))
	}
	finalWith := strings.Join(append(strings.Split(characters[:charUses], ""), "n"), ",")
	if len(matchStmt) > 0 {
		final := []string{
			"MATCH",
			strings.Join(matchStmt, ","),
			fmt.Sprintf(`CREATE (n:%s { %s })`, req.Type, strings.Join(vals, ",")),
			fmt.Sprintf(`WITH %s`, finalWith),
			"CREATE",
			strings.Join(createStmt, ","),
		}
		return strings.Join(final, " ")
	}
	return fmt.Sprintf(`CREATE (n:%s { %s })`, req.Type, strings.Join(vals, ","))

}
