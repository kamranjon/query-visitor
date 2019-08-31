package ast

// Resource is the root node in our inverted tree (all queries start here)
type Resource struct {
	Tables     []string
	Fields     []string
	ChildNodes []Node `json:"children"`
}

func (n *Resource) Accept(v Visitor) { v.Visit(n) }
func (n *Resource) Children() []Node { return n.ChildNodes }

type Filter struct {
	ChildNodes []Node
}

func (n *Filter) Accept(v Visitor) { v.Visit(n) }
func (n *Filter) Children() []Node { return n.ChildNodes }

type And struct {
	ChildNodes []Node
}

func (n *And) Accept(v Visitor) { v.Visit(n) }
func (n *And) Children() []Node { return n.ChildNodes }

type Or struct {
	ChildNodes []Node
}

func (n *Or) Accept(v Visitor) { v.Visit(n) }
func (n *Or) Children() []Node { return n.ChildNodes }

type GreaterThan struct {
	Property   string
	Value      interface{}
	ChildNodes []Node
}

func (n *GreaterThan) Accept(v Visitor) { v.Visit(n) }
func (n *GreaterThan) Children() []Node { return n.ChildNodes }

type LessThan struct {
	Property   string
	Value      interface{}
	ChildNodes []Node
}

func (n *LessThan) Accept(v Visitor) { v.Visit(n) }
func (n *LessThan) Children() []Node { return n.ChildNodes }

type EqualTo struct {
	Property   string
	Value      interface{}
	ChildNodes []Node
}

func (n *EqualTo) Accept(v Visitor) { v.Visit(n) }
func (n *EqualTo) Children() []Node { return n.ChildNodes }

type NotEqualTo struct {
	Property   string
	Value      interface{}
	ChildNodes []Node
}

func (n *NotEqualTo) Accept(v Visitor) { v.Visit(n) }
func (n *NotEqualTo) Children() []Node { return n.ChildNodes }
