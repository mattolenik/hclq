package query

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/parser"
)

type Result struct {
	Serialized string
	Node       ast.Node
}

func HCL(reader io.Reader, qry []Node) (results []Result, isList bool, err error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	node, err := parser.Parse(bytes)
	if err != nil {
		return
	}
	err = Walk(node.Node, qry, 0, func(n ast.Node, queryNode Node) (err error) {
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

func Walk(astNode ast.Node, query []Node, queryIdx int, action func(node ast.Node, queryNode Node) error) error {
	switch node := astNode.(type) {
	case *ast.ObjectList:
		for _, obj := range node.Items {
			err := Walk(obj, query, queryIdx, action)
			if err != nil {
				return err
			}
		}
		return nil

	case *ast.ObjectItem:
		queryLen := len(query)
		for _, key := range node.Keys {
			if !query[queryIdx].IsMatch(key.Token.Text, node.Val) {
				return nil
			}
			if queryIdx+1 >= queryLen {
				break
			}
			queryIdx++
		}
		// Assume a match if the for loop didn't return
		// Assume Keys will always be len > 0
		return Walk(node.Val, query, queryIdx, action)

	case *ast.ListType:
		return action(node, query[queryIdx])

	case *ast.LiteralType:
		return action(node, query[queryIdx])

	case *ast.ObjectType:
		return Walk(node.List, query, queryIdx, action)

	default:
		return errors.New("unhandled case")
	}
}
