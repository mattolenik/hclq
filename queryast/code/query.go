package code

type Query struct {
	Parts *Expr
}

type Expr struct {
	Node interface{}
	Next interface{}
}

type FunctionCall struct {
	Name   string
	Params interface{}
}

type Path struct {
	Crumbs []*Crumb
}

type Crumb struct {
	Key *Key
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
