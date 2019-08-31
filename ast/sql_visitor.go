package ast

import (
	"github.com/kamranjon/query-visitor/builder"
)

// SQLVisitor is the SQL implementation of the Visitor interface
type SQLVisitor struct {
	Query *builder.SQLBuilder
}

// Visit recursively calls node.Accept on the tree
func (sql *SQLVisitor) Visit(node Node) Visitor {
	switch node := node.(type) {
	case *Resource:
		sql.Resource(node)
	case *GreaterThan:
		sql.GreaterThan(node)
	case *LessThan:
		sql.LessThan(node)
	case *EqualTo:
		sql.EqualTo(node)
	case *And:

	}

	for _, c := range node.Children() {
		c.Accept(sql)
	}

	return nil
}

// Resource is the SQL specialization of Resource Node
func (sql *SQLVisitor) Resource(n *Resource) {
	sql.Query.Select(n.Fields...).From(n.Tables...)
}

// GreaterThan is the SQL specialization of GreaterThan Node
func (sql *SQLVisitor) GreaterThan(n *GreaterThan) {
	sql.Query.Where(sql.Query.Condition(n.Property).Greater(n.Value))
}

// LessThan is the SQL specialization of LessThan Node
func (sql *SQLVisitor) LessThan(n *LessThan) {
	sql.Query.Where(sql.Query.Condition(n.Property).Less(n.Value))
}

// EqualTo is the SQL specialization of EqualTo Node
func (sql *SQLVisitor) EqualTo(n *EqualTo) {
	sql.Query.Where(sql.Query.Condition(n.Property).Equal(n.Value))
}
