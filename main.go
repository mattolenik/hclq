package main

import (
	"fmt"
	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/hashicorp/hcl/hcl/ast"
	"io/ioutil"
	"github.com/docopt/docopt-go"
	"regexp"
	"os"
)

func main() {
	usage := `
HCL Query/Editor

Usage:
  hclq get <path> <file>
  hclq set <path>
  hclq --help
  hclq --version

Options:
  --help     Show this screen.
  --version  Show version.
`
	arguments, _ := docopt.Parse(usage, nil, true, "0.1.0-DEV", false)
	fmt.Println(arguments)

	var rt root
	get := arguments["get"].(bool)
	if get {
		rt = root {}
		parseQuery(arguments["<path>"].(string), 0, &rt)
	}
	fmt.Println(rt)
	bytes, err := ioutil.ReadFile(arguments["<file>"].(string)); check(err)
	node, err := parseAst(bytes); check(err)
	ast.Walk(node, func(node ast.Node) (ast.Node, bool) {
		//lt, ok := node.(*ast.LiteralType);
		//if ok {
		//	fmt.Println(lt.Token.Type, lt.Token.Value())
		//}
		//objectKey, ok := node.(*ast.ObjectKey)
		//if ok {
		//	fmt.Println(objectKey.GoString())
		//}
		objectItem, ok := node.(*ast.ObjectItem)
		if ok {
			fmt.Printf("%s %#v\n", objectItem.Val.Pos(), objectItem.Keys)
		}
		//objectList, ok := node.(*ast.ObjectList)
		//if ok {
		//	fmt.Println(objectList.Items)
		//}
		return node, true
	})
}

func check(err error) {
	if err != nil {
		fmt.Println(fmt.Errorf("%s", err))
		os.Exit(1)
	}
}

func parseAst(bytes []byte) (*ast.File, error) {
	node, err := parser.Parse(bytes);
	if err != nil {
		return nil, err
	}
	return node, nil
}

type queryNode interface {
	next() queryNode
	setNext(node queryNode)
	value() string
	setValue(value string)
	token() string
	setToken(value string)
}

func consume(s string) (node queryNode, err error) {
	if err != nil {
		return
	}
	tok := LiteralRegex.FindString(s)
	res := literal{
		tok: tok,
	}
	node = &res
	return
}

var LiteralRegex, _ = regexp.Compile(`(\w+)`)

// .simple.struct.val
// .foo.bar.items[*]
func parseQuery(query string, i int, node queryNode) {
	next := i + 1
	if next > len(query) {
		return
	}
	c := query[i:next]
	if c == "." {
		next += 1
		newNode := &dot {
			tokenStr: "dot",
		}
		node.setNext(newNode)
		parseQuery(query, next, newNode)
		return
	}
	word := LiteralRegex.FindString(query)
	if word != "" {
		newNode := &literal {
			tok: "literal",
			val: word,
		}
		node.setNext(newNode)
		next += len(word)
		parseQuery(query, next, newNode)
		return
	}
}