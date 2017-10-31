package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/hashicorp/hcl/hcl/token"
	jsonParser "github.com/hashicorp/hcl/json/parser"
	"github.com/mattolenik/hclq/query"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Automatically replaced when building with -ldflags="-X main.version=<version>"
var version = "undefined"

type ErrorJSON struct {
	Error string `json:"error"`
}

func main() {
	usage := `
HCL Query/Edit Tool

Usage:
  hclq get     [options] <node> <file>
  hclq get     [options] <node> [-]
  hclq set     [options] <node> <value> <file>
  hclq set     [options] <node> <value> [-]
  hclq append  [options] <node> <value> <file>
  hclq append  [options] <node> <value> [-]
  hclq prepend [options] <node> <value> <file>
  hclq prepend [options] <node> <value> [-]
  hclq replace [options] <node> <value> <newValue> <file>
  hclq replace [options] <node> <value> <newValue> [-]

  hclq --help
  hclq --version

Options:
  -i --in-place        Modify file in-place instead of writing to stdout
  -r --raw             Output in raw mode (Go printf %+v) instead of JSON
  -q --quiet           Ignore failures, output matches or nothing at all
  --help               Show this screen
  --version            Show version
`
	arguments, _ := docopt.Parse(usage, nil, true, version, false)
	queryNodes := query.Parse(arguments["<node>"].(string))

	raw := arguments["--raw"].(bool)
	inPlace := arguments["--in-place"].(bool)
	file, _ := arguments["<file>"].(string)
	value, valueOk := arguments["<value>"].(string)
	newValue, newValueOk := arguments["<newValue>"].(string)

	err := func() error {
		if arguments["get"].(bool) {
			return get(arguments, queryNodes, raw)

		} else if arguments["set"].(bool) {
			if !valueOk {
				return errors.New("<value> required for set command")
			}
			return set(
				file,
				queryNodes,
				func(s string) string { return value },
				inPlace)

		} else if arguments["append"].(bool) {
			if !valueOk {
				return errors.New("<value> required for append command")
			}
			return set(
				file,
				queryNodes,
				func(s string) string { return s + value },
				inPlace)

		} else if arguments["prepend"].(bool) {
			if !valueOk {
				return errors.New("<value> required for prepend command")
			}
			return set(
				file,
				queryNodes,
				func(s string) string { return value + s },
				inPlace)

		} else if arguments["replace"].(bool) {
			if !valueOk {
				return errors.New("<value> required for replace command")
			}
			if !newValueOk {
				return errors.New("<newValue> required for replace command")
			}
			return set(
				file,
				queryNodes,
				func(s string) string {
					return strings.Replace(s, value, newValue, -1)
				},
				inPlace)
		}
		return errors.New("unknown command")
	}()
	if err != nil && !arguments["--quiet"].(bool) {
		getErrorOutput(err, raw)
		os.Exit(1)
	}
}

func getErrorOutput(err error, raw bool) string {
	if raw {
		return fmt.Sprintf("%+v", err)
	} else {
		out, _ := json.Marshal(ErrorJSON{Error: err.Error()})
		return string(out)
	}
}

func getOutput(obj interface{}, raw bool) (string, error) {
	if raw {
		return fmt.Sprintf("%+v", obj), nil
	} else {
		jsonBody, err := json.Marshal(obj)
		if err != nil {
			return "", errors.New("failure while trying to serialize output to JSON")
		}
		return string(jsonBody), nil
	}
}

func get(arguments map[string]interface{}, query []query.Node, raw bool) error {
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

	output, err := getOutput(result, raw)
	if err != nil {
		return err
	}

	fmt.Print(output)
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
		if len(result) == 0 {
			return nil, errors.New("no match found for query")
		}
		if len(result) == 1 {
			return result[0], nil
		}
		return result, nil
	} else if objItem, ok := node.(*ast.ObjectItem); ok {
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
	} else if objType, ok := node.(*ast.ObjectType); ok {
		return getImpl(objType.List, query, queryIdx)
	} else {
		return toGoType(node)
	}
	return nil, errors.New("unhandled case")
}

func set(fileName string, query []query.Node, value func(original string) string, inPlace bool) error {
	var hcl []byte
	var err error

	if fileName != "" {
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

	err = setImpl(node.Node, query, value, 0)
	if err != nil {
		return err
	}

	if inPlace {
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

func setImpl(node ast.Node, query []query.Node, value func(original string) string, queryIdx int) error {
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
		literal.Token.Text = value(literal.Token.Text)
		return nil
	}
	if list, ok := node.(*ast.ListType); ok {
		// HCL JSON parser needs a top level object
		jsonValue := fmt.Sprintf(`{"root": %s}`, )
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
	return errors.New("unhandled case")
}

func toGoType(node ast.Node) (interface{}, error) {
	if literal, ok := node.(*ast.LiteralType); ok {
		switch literal.Token.Type {
		case token.STRING:
			return literal.Token.Value().(string), nil
		case token.HEREDOC:
			return literal.Token.Value().(string), nil
		case token.FLOAT:
			return literal.Token.Value().(float64), nil
		case token.NUMBER:
			return literal.Token.Value().(int64), nil
		case token.BOOL:
			return literal.Token.Value().(bool), nil
		}
	} else if list, ok := node.(*ast.ListType); ok {
		var result []interface{}
		for _, item := range list.List {
			nextItem, err := toGoType(item)
			if err != nil {
				return nil, err
			}
			result = append(result, nextItem)
		}
		return result, nil
	} else if objectItem, ok := node.(*ast.ObjectItem); ok {
		result, err := toGoType(objectItem.Val)
		return result, err
	}
	return "", errors.New("unhandled type conversion")
}
