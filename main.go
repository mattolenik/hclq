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
		FilterSeparator     = "|" .
		QuerySeparator      = "." | Whitespace "." .
		Regex               = "/" { "\\/" | anyregex } "/" .
		Ident               = { "\\." | anyident } | { anyident } .
		Whitespace          = "\r\n" | "\n" | "\r" | " " | "\t" .

		alpha           = "a"…"z" | "A"…"Z" .
		invalid         = "\u0000"…"\u0040" | "\u007B"…"\u00BF" .
		any             = "\u0041"…"\u007A" | "\u00C0"…"\uFFFF" .
		anyident        = "\u0041"…"\u007A" - "." | "\u00C0"…"\uFFFF" .
		anyregex        = "\u0041"…"\u007A" - "/" | "\u00C0"…"\uFFFF" .
	`))

	parser, err := participle.Build(
		&Query{},
		participle.Lexer(lxr),
		participle.Elide("Whitespace", "invalid"),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	q := `.abc
.def
./abc\/sadf/.xy\.z`
	qAST := &Query{}
	err = parser.ParseString(q, qAST)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	spew.Dump(qAST)
}

type Query struct {
	Crumbs []*Crumb `( QuerySeparator @@ )+`
}

type Crumb struct {
	Ident string `  @Ident`
	Regex string `| @Regex`
}
