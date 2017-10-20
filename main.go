package main

import (
	"fmt"
	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/hashicorp/hcl/hcl/ast"
	"io/ioutil"
	"github.com/docopt/docopt-go"
	"regexp"
	"os"
	"strings"
	"github.com/davecgh/go-spew/spew"
	"reflect"
	"strconv"
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

	get := arguments["get"].(bool)
	query := make([]QueryNode, 0)
	if get {
		parseQuery(arguments["<path>"].(string), 0, &query)
	}

	bytes, err := ioutil.ReadFile(arguments["<file>"].(string));
	check(err)

	node, err := parseAst(bytes);
	check(err)

	//spew.Dump(node)

	var result interface{}
	done := false
	errResult := "Unable to satisfy query"
	idxStack := []int{0}

	ast.Walk(node, func(node ast.Node) (ast.Node, bool) {
		if done {
			return node, false
		}

		stackLen := len(idxStack)
		qi := idxStack[stackLen - 1]
		if stackLen > 0 {
			idxStack = idxStack[:stackLen - 1]
		}

		objectItem, ok := node.(*ast.ObjectItem)
		if ok {
			for _, key := range objectItem.Keys {
				value := strings.Trim(key.Token.Text, "\"")
				if qi >= len(query) || query[qi].Value() != value {
					return node, false
				}
				qi++
			}
			idxStack = append(idxStack, qi + 1)
			if qi >= len(query) {
				result, err = toGoType(node)
				if err != nil {
					errResult = err.Error()
					done = true
					return node, false
				}
			}

			return node, true
		}
		//	objectList, ok := node.(*ast.ObjectList)
		//	if ok {
		//		spew.Dump(objectList.Items[0].Val)
		//	}
		return node, true
	})

	if result == nil {
		fmt.Fprintln(os.Stderr, errResult)
	} else {
		fmt.Printf("Type: %s, Value: %v", reflect.TypeOf(result), result)
	}
}

func toGoType(node ast.Node) (interface{}, error) {
	if literal, ok := node.(*ast.LiteralType); ok {
		token := literal.Token.Text
		num, err := strconv.ParseUint(token, 10, 64)
		if err == nil {
			return num, nil
		}
		return token, nil
	}
	if list, ok := node.(*ast.ListType); ok {
		var result []interface{}
		for _, item := range list.List {
			nextItem, err := toGoType(item)
			if err != nil {
				return nil, err
			}
			result = append(result, nextItem)
		}
		return result, nil
	}
	if objectItem, ok := node.(*ast.ObjectItem); ok {
		result, err := toGoType(objectItem.Val)
		return result, err
	}
	spew.Dump(node)
	return "", nil
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

func consume(s string) (node QueryNode, err error) {
	if err != nil {
		return
	}
	tok := LiteralRegex.FindString(s)
	res := Key{
		value: tok,
	}
	node = &res
	return
}

var LiteralRegex, _ = regexp.Compile(`(\w+)`)

// .simple.struct.value
// .foo.bar.items[*]
func parseQuery(query string, i int, queue *[]QueryNode) {
	if i >= len(query) {
		return
	}
	char := query[i:i+1]
	if char == "." {
		parseQuery(query, i+1, queue)
		return
	}
	word := LiteralRegex.FindString(query[i:])
	if word != "" {
		i += len(word)
		newNode := &Key{
			value:      word,
		}
		*queue = append(*queue, newNode)
		if i >= len(query) {
			return
		}
		parseQuery(query, i, queue)
		return
	}
}
