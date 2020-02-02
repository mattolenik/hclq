package code

import "github.com/hashicorp/hcl2/hcl/hclsyntax"

type Query struct {
	Parts *Expr
}

type Expr struct {
	Node interface{}
	Next interface{}
}

type Exp interface {
	Input(node hclsyntax.Node)
	Evaluate()
}

type GatherExpr interface {
	Gather() interface{}
	Node() interface{}
}

type ResultsExpr interface {
	Results() interface{}
}

type ExecuteExpr interface {
	SetInput(GatherExpr)
	Execute() ResultsExpr
}

type nextExpr struct {
	Separator interface{}
	Expr      *Expr
}

type PipeOperator struct {
}

type DescentOperator struct {
}

type FunctionCall struct {
	Name   string
	Params interface{}
}

type Crumb struct {
	Key  *Key
	Next interface{}
}

type EmptySelector struct {
}

func (s *EmptySelector) Accumulate(accum *Accumulation) error {
	return nil
}

type IndexSelector struct {
	Index int
}

type SplatSelector struct {
}

type MapSelector struct {
	Key string
}

type Key struct {
	Ident    string
	Selector interface{}
}

type Accumulation struct {
	Results []interface{}
}

func (a *Accumulation) AddResult(result interface{}) {
	a.Results = append(a.Results, result)
}

type Accumulator interface {
	AddResult() func(a *Accumulator) error
}
