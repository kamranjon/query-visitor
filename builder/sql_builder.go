package builder

import (
	"fmt"
	"strings"
)

type SQLBuilder struct {
	Columns       []string
	Tables        []string
	WhereClause   *WhereClause
	GroupByClause map[string]string
	JoinClauses   []JoinClause
	Parameters    []interface{}
}

type WhereClause struct {
	LogicalOperator string
	Conditions      []*Condition
	SubClause       *WhereClause
}

type JoinClause struct {
	Table     string
	Condition Condition
}

type Condition struct {
	Property         string
	Value            interface{}
	Operator         string
	LogicalOperator  string
	SubConditionNode *Condition
}

func (c *Condition) applyOperator(value []interface{}, operator string) *Condition {
	sc := c.SubConditionNode
	if sc != nil {
		sc.applyOperator(value, operator)
	} else {
		c.Operator = operator
		if len(value) != 0 {
			c.Value = value[0]
		}
	}

	return c
}

func lastCondition(c *Condition) *Condition {
	if c.SubConditionNode == nil {
		return c
	}
	return lastCondition(c.SubConditionNode)
}

func (s *SQLBuilder) LastCondition() *Condition {
	lastWhere := s.LastWhere()
	conditions := lastWhere.Conditions
	return lastCondition(conditions[len(conditions)-1])
}

// Conditional Operators
func (c *Condition) Greater(value ...interface{}) *Condition {
	return c.applyOperator(value, ConditionalOperator.Greater)
}

func (c *Condition) Less(value ...interface{}) *Condition {
	return c.applyOperator(value, ConditionalOperator.Less)
}

func (c *Condition) GreaterOrEqual(value ...interface{}) *Condition {
	return c.applyOperator(value, ConditionalOperator.GreaterOrEqual)
}

func (c *Condition) LessOrEqual(value ...interface{}) *Condition {
	return c.applyOperator(value, ConditionalOperator.LessOrEqual)
}

func (c *Condition) Equal(value ...interface{}) *Condition {
	return c.applyOperator(value, ConditionalOperator.Equal)
}

func (c *Condition) NotEqual(value ...interface{}) *Condition {
	return c.applyOperator(value, ConditionalOperator.NotEqual)
}

// Arithmetic Operators
func (c *Condition) Plus(value ...interface{}) *Condition {
	return c.applyOperator(value, ArithmeticOperator.Plus)
}

func (c *Condition) Minus(value ...interface{}) *Condition {
	return c.applyOperator(value, ArithmeticOperator.Minus)
}

func (c *Condition) Divide(value ...interface{}) *Condition {
	return c.applyOperator(value, ArithmeticOperator.Divide)
}

func (c *Condition) Multiply(value ...interface{}) *Condition {
	return c.applyOperator(value, ArithmeticOperator.Multiply)
}

func (c *Condition) Modulo(value ...interface{}) *Condition {
	return c.applyOperator(value, ArithmeticOperator.Modulo)
}

func (c *Condition) Condition(property string) *Condition {
	sc := lastCondition(c)
	sc.SubConditionNode = &Condition{Property: property}
	return c
}

func (s *SQLBuilder) Condition(property string) *Condition {
	return &Condition{Property: property}
}

var ColumnOrder = struct {
	Asc  string
	Desc string
}{
	"ASC",
	"DESC",
}

var ArithmeticOperator = struct {
	Plus     string
	Minus    string
	Multiply string
	Divide   string
	Modulo   string
}{
	"+",
	"-",
	"*",
	"/",
	"%",
}

var ConditionalOperator = struct {
	Greater        string
	Less           string
	Equal          string
	GreaterOrEqual string
	LessOrEqual    string
	NotEqual       string
}{
	">",
	"<",
	"=",
	">=",
	"<=",
	"<>",
}

var LogicalOperator = struct {
	And string
	Or  string
	Not string
}{
	" AND ",
	" OR ",
	"!",
}

func NewQuery() *SQLBuilder {
	return new(SQLBuilder)
}

func (s *SQLBuilder) Select(columns ...string) *SQLBuilder {
	s.Columns = columns
	return s
}

func (s *SQLBuilder) From(tables ...string) *SQLBuilder {
	s.Tables = tables
	return s
}

func (s *SQLBuilder) Where(conditions ...*Condition) *SQLBuilder {
	// if no clauses exist - start a new where clause
	if s.WhereClause == nil {
		s.WhereClause = &WhereClause{
			LogicalOperator: LogicalOperator.And,
			Conditions:      conditions,
			SubClause:       nil,
		}
		return s
	}

	lastWhereClause := s.LastWhere()
	// otherwise - append to the already in-progress where clause
	lastWhereClause.SubClause = &WhereClause{
		LogicalOperator: LogicalOperator.And,
		Conditions:      conditions,
		SubClause:       nil,
	}
	return s
}

func lastWhere(where *WhereClause) *WhereClause {
	if where.SubClause == nil {
		return where
	}
	return lastWhere(where.SubClause)
}

func (s *SQLBuilder) LastWhere() *WhereClause {
	return lastWhere(s.WhereClause)
}

func (s *SQLBuilder) And(property string) *Condition {
	return &Condition{Property: property, LogicalOperator: LogicalOperator.And}
}

func (s *SQLBuilder) Or(property string) *Condition {
	return &Condition{Property: property, LogicalOperator: LogicalOperator.Or}
}

func (s *SQLBuilder) Not(property string) *Condition {
	return &Condition{Property: property, LogicalOperator: LogicalOperator.Not}
}

func (s *SQLBuilder) Join(table string, property string, operator string, value interface{}) *SQLBuilder {
	condition := Condition{
		Property: property,
		Value:    value,
		Operator: operator,
	}
	s.JoinClauses = append(s.JoinClauses, JoinClause{Table: table, Condition: condition})
	return s
}

func (s *SQLBuilder) GroupBy(columnOrders ...string) *SQLBuilder {
	if len(columnOrders) == 1 {
		s.GroupByClause = map[string]string{columnOrders[0]: ColumnOrder.Asc}
		return s
	}

	newMap := map[string]string{}
	for i := range columnOrders {
		if i+1%2 == 0 {
			newMap[columnOrders[i-1]] = columnOrders[i]
		}
	}
	s.GroupByClause = newMap
	return s
}

func (s *SQLBuilder) buildSelect(columns []string) string {
	return fmt.Sprintf("SELECT %s", strings.Join(columns, ", "))
}

func (s *SQLBuilder) buildFrom(tables []string) string {
	return fmt.Sprintf("FROM %s", strings.Join(tables, ", "))
}

func (s *SQLBuilder) buildCondition(condition *Condition, strings []string) []string {
	if condition.Value != nil {
		s.Parameters = append(s.Parameters, condition.Value)
		conditionString := fmt.Sprintf(
			"%s(%s %s ?)",
			condition.LogicalOperator,
			condition.Property,
			condition.Operator,
		)
		strings = append(strings, conditionString)
		return strings
	}

	if condition.SubConditionNode != nil {
		conditionString := fmt.Sprintf("(%s %s ", condition.Property, condition.Operator)
		strings = append(strings, conditionString)
		strings = append(strings, s.buildCondition(condition.SubConditionNode, []string{})...)
		strings = append(strings, ")")
	}

	return strings

}

func (s *SQLBuilder) buildWhere(where *WhereClause) string {
	conditions := where.Conditions
	conditionStrings := []string{}

	for _, condition := range conditions {
		conditionList := s.buildCondition(condition, []string{})
		conditionString := strings.Join(conditionList, "")
		conditionStrings = append(conditionStrings, conditionString)
	}

	if len(conditions) > 1 {
		conditionStrings = append([]string{"("}, conditionStrings...)
		conditionStrings = append(conditionStrings, ")")
	}

	if where.SubClause != nil {
		subClauseString := fmt.Sprintf("%s%s", LogicalOperator.And, s.buildWhere(where.SubClause))
		conditionStrings = append(conditionStrings, subClauseString)
	}

	return strings.Join(conditionStrings, "")
}

func (s *SQLBuilder) ToSQL() (string, []interface{}) {
	selectString := s.buildSelect(s.Columns)
	fromString := s.buildFrom(s.Tables)
	whereString := s.buildWhere(s.WhereClause)

	queryString := fmt.Sprintf("%s %s WHERE %s", selectString, fromString, whereString)
	return queryString, s.Parameters
}
