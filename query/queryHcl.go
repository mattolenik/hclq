package query

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/parser"
)

// Result represents a query result
type Result struct {
	Value interface{}
	Node  ast.Node
}

// HCL performs a generic query and returns matching results
func HCL(reader io.Reader, qry *Query) (results []Result, isList bool, node *ast.File, err error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	node, err = parser.Parse(bytes)
	if err != nil {
		return
	}
	err = walk(node.Node, qry, func(n ast.Node, queryNode Node) (err error) {
		switch node := n.(type) {
		case *ast.LiteralType:
			switch queryNode.(type) {
			case Node:
			case IndexedNode:
				return fmt.Errorf("Invalid query, '%s', matching item is not an array", queryNode.Key())
			default:
				return fmt.Errorf("Query format '%s' not understood", queryNode.Key())
			}
			results = append(results, Result{node.Token.Value(), node})

		case *ast.ListType:
			listNode, ok := queryNode.(IndexedNode)
			if !ok {
				return fmt.Errorf("unexpected query type")
			}
			// Query is for a specific index
			if listNode.Index() != nil {
				listLength := len(node.List)
				listIndex := *listNode.Index()

				// Negative index means wrap around, with -1 being the last element
				if listIndex < 0 {
					listIndex = listLength + listIndex
				}
				if listIndex < 0 || listIndex >= listLength {
					return fmt.Errorf("index %d out of bounds on list %+v of len %d", listIndex, listNode.Key(), listLength)
				}
				val, ok := node.List[listIndex].(*ast.LiteralType)
				if !ok {
					return err
				}
				results = append(results, Result{val.Token.Value(), node})
				return nil
			}
			// Query is for all elements
			isList = true
			for _, item := range node.List {
				if literal, ok := item.(*ast.LiteralType); ok {
					results = append(results, Result{literal.Token.Value(), node})
				}
			}
		// TODO: full objects
		//case *ast.ObjectItem:
		default:
			fmt.Println(node)
		}
		return nil
	})
	return
}

func walk(astNode ast.Node, query *Query, action func(node ast.Node, queryNode Node) error) error {
	return walkImpl(astNode, query, 0, action)
}

func walkImpl(astNode ast.Node, query *Query, qIdx int, action func(node ast.Node, queryNode Node) error) error {
	switch node := astNode.(type) {
	case *ast.ObjectList:
		for _, obj := range node.Items {
			err := walkImpl(obj, query, qIdx, action)
			if err != nil {
				return err
			}
		}
		return nil

	case *ast.ObjectItem:
		for _, key := range node.Keys {
			isMatch := query.Parts[qIdx].IsMatch(strings.Trim(key.Token.Text, `"`), node.Val)
			if !isMatch {
				return nil
			}
			if qIdx+1 >= query.Length {
				break
			}
			qIdx++
		}
		// Assume a match if return didn't happen in the for loop.
		// Assume Keys will always be len > 0 (it wouldn't be valid HCL otherwise)
		return walkImpl(node.Val, query, qIdx, action)

	case *ast.ListType:
		return action(node, query.Parts[qIdx])

	case *ast.LiteralType:
		return action(node, query.Parts[qIdx])

	case *ast.ObjectType:
		return walkImpl(node.List, query, qIdx, action)

	default:
		return errors.New("unhandled case")
	}
}
