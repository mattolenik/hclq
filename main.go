package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	lxr := lexer.Must(ebnf.New(`
		Dot = "." .
		Regex = "/" { "\\/" | anyregex } "/" .
		Ident = { "\\." | anyident } | { anyident } .
		alpha = "a"…"z" | "A"…"Z" .
		any = "\u0000"…"\uffff" .
		anyident = "\u0000"…"\uffff"-"." .
		anyregex = "\u0000"…"\uffff"-"/" .
	`))

	parser, err := participle.Build(
		&Query{},
		participle.Lexer(lxr),
	)
	if err != nil {
		panic(err)
	}
	q := `.abc.def./abc\/sadf/.xy\.z`
	qAST := &Query{}
	err = parser.ParseString(q, qAST)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	spew.Dump(qAST)
}

type Query struct {
	Crumbs []*Crumb `( Dot @@ )+`
}

type Crumb struct {
	Ident string `  @Ident`
	Regex string `| @Regex`
}
