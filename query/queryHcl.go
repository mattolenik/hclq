package query

import (
	"io"
	"io/ioutil"
	"github.com/mattolenik/hclq/utils"
	"fmt"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/parser"
)

type Result struct {
	Serialized string
	Node ast.Node
}

func QueryHcl(reader io.Reader, qry []Node) (results []Result, isList bool, err error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	node, err := parser.Parse(bytes)
	if err != nil {
		return
	}
	err = utils.Walk(node.Node, qry, 0, func(n ast.Node, queryNode Node) (err error) {
		switch node := n.(type) {
		case *ast.LiteralType:
			results = append(results, Result{node.Token.Text, node})

		case *ast.ListType:
			listNode, ok := queryNode.(*List)
			if !ok {
				return fmt.Errorf("unexpected query type")
			}
			// Query is for a specific index
			if listNode.Index != nil {
				listLength := len(node.List)
				listIndex := *listNode.Index
				if listIndex >= listLength {
					return fmt.Errorf("index %d out of bounds on list %+v of len %d", listIndex, listNode.Key, listLength)
				}
				val, ok := node.List[listIndex].(*ast.LiteralType)
				if !ok {
					return err
				}
				results = append(results, Result{val.Token.Text, node})
				return nil
			}
			// Query is for all elements
			isList = true
			for _, item := range node.List {
				if literal, ok := item.(*ast.LiteralType); ok {
					results = append(results, Result{literal.Token.Text, node})
				}
			}
		default:
			fmt.Println(node)
		}
		return nil
	})
	return
}
