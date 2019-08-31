package main

import (
	//"encoding/json"
	"log"

	"github.com/kamranjon/query-visitor/ast"
	"github.com/kamranjon/query-visitor/builder"
)

func main() {
	tree :=
		&ast.Resource{
			Tables: []string{"messages"},
			Fields: []string{"*"},
			ChildNodes: []ast.Node{
				&ast.And{
					ChildNodes: []ast.Node{
						&ast.GreaterThan{
							Property: "score",
							Value:    12,
							ChildNodes: []ast.Node{
								&ast.LessThan{
									Property: "height",
									Value:    25,
								},
							},
						},
					},
				},
			},
		}

	s := builder.NewQuery()
	s.Select("*").
		From("people").
		Where(s.
			Condition("score").
			Greater().
			Condition("age").
			Plus().
			Condition("wealth").
			Minus().
			Condition("heart").
			Multiply(10),
			s.And("friends").
				Greater(10),
		).
		Where(s.
			Condition("gold").
			Less(25),
			s.Or("appetite").
				Less(12),
		).
		Where(s.
			Condition("old_friends").
			Equal(0),
		)

	sql, vars := s.ToSQL()

	log.Println("TO SQL:", sql, vars)
	visitor := ast.SQLVisitor{Query: builder.NewQuery()}

	visitor.Visit(tree)
	preparedQuery, params := visitor.Query.ToSQL()
	log.Println(preparedQuery)
	log.Println(params)
}
