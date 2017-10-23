package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/hashicorp/hcl/hcl/printer"
	jsonParser "github.com/hashicorp/hcl/json/parser"
	"github.com/mattolenik/hclq/query"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Automatically replaced when building with -ldflags="-X main.version=<version>"
var version = "undefined"

func main() {
	usage := `
HCL Query/Editor

Usage:
  hclq get <file> <nodePath>
  hclq get <nodePath>
  hclq set [-i] <file> <nodePath> <value>
  hclq set <nodePath> <value>
  hclq --help
  hclq --version

Options:
  -i        Modify file in-place instead of writing to stdout
  --help     Show this screen
  --version  Show version
`
	arguments, _ := docopt.Parse(usage, nil, true, version, true)
	queryNodes := query.Parse(arguments["<nodePath>"].(string))

	var err error
	if arguments["get"].(bool) {
		err = get(arguments, queryNodes)
	}
	if arguments["set"].(bool) {
		err = set(arguments, queryNodes)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}

func get(arguments map[string]interface{}, query []query.Node) error {
	var reader io.Reader
	fileName, ok := arguments["<file>"].(string)
	if !ok {
		reader = os.Stdin
	} else {
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		reader = bufio.NewReader(file)
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	node, err := parser.Parse(bytes)
	if err != nil {
		return err
	}

	result, err := getImpl(node.Node, query, 0)
	if err != nil {
		return err
	}

	fmt.Printf("%+v", result)
	return nil
}

func getImpl(node ast.Node, query []query.Node, queryIdx int) (interface{}, error) {
	if objList, ok := node.(*ast.ObjectList); ok {
		var result []interface{}
		for _, obj := range objList.Items {
			res, err := getImpl(obj, query, queryIdx)
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
			value := strings.Trim(key.Token.Text, "\"")
			queryKey := query[queryIdx].Value()
			if value != queryKey {
				return nil, nil
			}
			queryIdx++
		}
		// Assume a match if the for loop didn't return
		// Assume Keys will always be len > 0
		return getImpl(objItem.Val, query, queryIdx)
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
		return getImpl(objType.List, query, queryIdx)
	}
	return nil, errors.New("Unhandled case")
}

func set(arguments map[string]interface{}, query []query.Node) error {
	var hcl []byte
	var err error
	var fileName string
	var ok bool

	if fileName, ok = arguments["<file>"].(string); ok {
		hcl, err = ioutil.ReadFile(fileName)
		if err != nil {
			return err
		}
	} else {
		hcl, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
	}

	node, err := parser.Parse(hcl)
	if err != nil {
		return err
	}

	err = setImpl(node.Node, query, arguments["<value>"].(string), 0)
	if err != nil {
		return err
	}

	if arguments["-i"].(bool) {
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		printer.Fprint(file, node)
	} else {
		printer.Fprint(os.Stdout, node)
	}
	return nil
}

func setImpl(node ast.Node, query []query.Node, value string, queryIdx int) error {
	if objList, ok := node.(*ast.ObjectList); ok {
		for _, obj := range objList.Items {
			err := setImpl(obj, query, value, queryIdx)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if objItem, ok := node.(*ast.ObjectItem); ok {
		queryLen := len(query)
		for _, key := range objItem.Keys {
			if queryIdx >= queryLen {
				return nil
			}
			value := strings.Trim(key.Token.Text, "\"")
			queryKey := query[queryIdx].Value()
			if value != queryKey {
				return nil
			}
			queryIdx++
		}
		// Assume a match if the for loop didn't return
		// Assume Keys will always be len > 0
		return setImpl(objItem.Val, query, value, queryIdx)
	}
	if literal, ok := node.(*ast.LiteralType); ok {
		literal.Token.Text = value
		return nil
	}
	if list, ok := node.(*ast.ListType); ok {
		// HCL JSON parser needs a top level object
		jsonValue := fmt.Sprintf(`{"root": %s}`, value)
		tree, err := jsonParser.Parse([]byte(jsonValue))
		if err != nil {
			return err
		}
		list.List = tree.Node.(*ast.ObjectList).Items[0].Val.(*ast.ListType).List
		return nil
	}
	if objType, ok := node.(*ast.ObjectType); ok {
		return setImpl(objType.List, query, value, queryIdx)
	}
	return errors.New("Unhandled case")
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
	return "", nil
}
