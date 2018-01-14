package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"errors"
	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/mattolenik/hclq/query"
	"github.com/mattolenik/hclq/utils"
	"io"
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/hashicorp/hcl/hcl/ast"
)

var GetCmd = &cobra.Command{
	Use:   "get <query>",
	Short: "retrieve matching values",
	RunE: func(cmd *cobra.Command, args []string) error {
		queryNodes, _ := query.Parse(args[0])
		reader := os.Stdin
		if val := cmd.Flag("in").Value.String(); val != "" {
			var err error
			reader, err = os.Open(val)
			if err != nil {
				return err
			}
		}
		err := get(reader, queryNodes, false)
		if err != nil {
			return err
		}
		return nil
	},
}

func getErrorOutput(err error, raw bool) string {
	if raw {
		return fmt.Sprintf("%+v", err)
	} else {
		out, _ := json.Marshal(utils.ErrorJSON{Error: err.Error()})
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

func get(reader io.Reader, qry []query.Node, raw bool) error {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	node, err := parser.Parse(bytes)
	if err != nil {
		return err
	}
	var results []string
	isList := false
	_, err = utils.Walk(node.Node, qry, 0, func(n ast.Node, queryNode query.Node) (stop bool, err error) {
		switch node := n.(type) {

		case *ast.LiteralType:
			results = append(results, node.Token.Text)

		case *ast.ListType:
			listNode, ok := queryNode.(*query.List)
			if !ok {
				return false, errors.New("unexpected query type")
			}
			// Query is for a specific index
			if listNode.Index != nil {
				listLength := len(node.List)
				listIndex := *listNode.Index
				if listIndex >= listLength {
					return true, fmt.Errorf("index %d out of bounds on list %s of len %d", listNode.Value, listIndex, listLength)
				}
				val, ok := node.List[listIndex].(*ast.LiteralType)
				if !ok {
					return false, err
				}
				results = append(results, val.Token.Text)
				return false, nil
			}
			// Query is for all elements
			isList = true
			for _, item := range node.List {
				if literal, ok := item.(*ast.LiteralType); ok {
					results = append(results, literal.Token.Text)
				}
			}
		default:
			fmt.Println(node)
		}
		return false, nil
	})
	if err != nil {
		return err
	}

	// The return type can be a list if: the queried object IS a list, or if the query matched multiple single items
	// So, return now if it's not a list and there is only one query result
	if !isList && len(results) == 1 {
		output, err := getOutput(results[0], raw)
		if err != nil {
			return err
		}
		fmt.Print(output)
		return nil
	}
	output, err := getOutput(results, raw)
	if err != nil {
		return err
	}
	fmt.Print(output)
	return nil
}