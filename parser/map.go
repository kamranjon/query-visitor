package parser

import "github.com/kamranjon/query-visitor/ast"

type NodeMapper struct{}

func (*NodeMapper) Get(sh string) ast.Node {
	switch sh {
	case "gt":
		return &ast.GreaterThan{}
	case "gte":
		return &ast.GreaterThanOrEqual{}
	case "lte":
		return &ast.LessThanOrEqual{}
	case "lt":
		return &ast.LessThan{}
	case "eq":
		return &ast.EqualTo{}
	case "and":
		return &ast.And{}
	case "or":
		return &ast.Or{}
	case "filter":
		return &ast.Filter{}
	default:
		return &ast.Resource{}
	}
}
