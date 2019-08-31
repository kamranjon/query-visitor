package ast

// Visitor defines the interface for our Visitor pattern
type Visitor interface {
	Visit(node Node) (w Visitor)
}
