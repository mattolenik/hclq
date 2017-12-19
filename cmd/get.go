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
	"container/list"
)

var GetCmd = &cobra.Command{
	Use:   "get <query>",
	Short: "retrieve matching values",
	Run: func(cmd *cobra.Command, args []string) {
		queryNodes, _ := query.Parse(args[0])
		reader := os.Stdin
		if val := cmd.Flag("in").Value.String(); val != "" {
			var err error
			reader, err = os.Open(val)
			if err != nil {
				panic(err)
			}
		}
		err := get(reader, queryNodes, false)
		if err != nil {
			panic(err)
		}
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

func get(reader io.Reader, query []query.Node, raw bool) error {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	node, err := parser.Parse(bytes)
	if err != nil {
		return err
	}
	literals := list.New()
	lists := list.New()
	_, err = utils.Walk(node.Node, query, 0, func(n ast.Node) (stop bool, err error) {
		if lit, ok := n.(*ast.LiteralType); ok {
			literals.PushBack(lit.Token.Text)
		}
		if lst, ok := n.(*ast.ListType); ok {
			lists.PushBack(lst.List)
		}
		return false, nil
	})
	if err != nil {
		return err
	}

	for literal := literals.Front(); literal != nil; literal = literal.Next() {
		output, err := getOutput(literal.Value, raw)
		if err != nil {
			return err
		}
		fmt.Print(output)
	}
	return nil
}