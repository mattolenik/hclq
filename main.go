package main

import (
	"fmt"
	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/hashicorp/hcl/hcl/ast"
	"io/ioutil"
	"github.com/docopt/docopt-go"
	"regexp"
	"os"
	"github.com/davecgh/go-spew/spew"
	"strconv"
	"strings"
	"errors"
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

	if arguments["get"].(bool) {
		query := make([]QueryNode, 0)
		parseQuery(arguments["<path>"].(string), 0, &query)
		bytes, err := ioutil.ReadFile(arguments["<file>"].(string));
		check(err)

		node, err := parseAst(bytes);
		check(err)

		result, err := get(node.Node, query, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("%+v", result)
	}
}

func get(node ast.Node, query []QueryNode, queryIdx int) (interface{}, error) {
	if objList, ok := node.(*ast.ObjectList); ok {
		var result []interface{}
		for _, obj := range objList.Items {
			res, err := get(obj, query, queryIdx)
			if err != nil {
				return nil, err
			}
			if res != nil {
				result = append(result, res)
			}
		}
		if len(result) == 1 {
			return result[0], nil
		}
		return result, nil
	}
	if objItem, ok := node.(*ast.ObjectItem); ok {
		queryLen := len(query)
		for _, key := range objItem.Keys {
			if queryIdx >= queryLen {
				return nil, nil
			}
			value := strings.Trim(key.Token.Text,"\"")
			queryKey := query[queryIdx].Value()
			if value != queryKey {
				return nil, nil
			}
			queryIdx++
		}
		// Assume a match if the for loop didn't return
		// Assume Keys will always be len > 0
		return get(objItem.Val, query, queryIdx)
	}
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
	if objType, ok := node.(*ast.ObjectType); ok {
		return get(objType.List, query, queryIdx)
	}
	return nil, errors.New("Unhandled case")
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
