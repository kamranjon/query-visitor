package ast

// Node defines the interface for our AST nodes
type Node interface {
	Accept(v Visitor)
	Children() []Node
}
