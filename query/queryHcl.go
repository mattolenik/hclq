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

type Results struct {
	Values []interface{}
	Node   ast.Node
}

func HCL(reader io.Reader, qry *Query) (results Results, isList bool, node *ast.File, err error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	node, err = parser.Parse(bytes)
	if err != nil {
		return
	}
	err = Walk(node.Node, qry, func(n ast.Node, queryNode Node) (err error) {
		switch node := n.(type) {
		case *ast.LiteralType:
			results.Node = node
			results.Values = append(results.Values, node.Token.Value())

		case *ast.ListType:
			listNode, ok := queryNode.(IndexedNode)
			if !ok {
				return fmt.Errorf("unexpected query type")
			}
			// Query is for a specific index
			if listNode.Index() != nil {
				listLength := len(node.List)
				listIndex := *listNode.Index()
				if listIndex >= listLength {
					return fmt.Errorf("index %d out of bounds on list %+v of len %d", listIndex, listNode.Key(), listLength)
				}
				val, ok := node.List[listIndex].(*ast.LiteralType)
				if !ok {
					return err
				}
				results.Node = node
				results.Values = append(results.Values, val.Token.Value())
				return nil
			}
			// Query is for all elements
			isList = true
			results.Node = node
			for _, item := range node.List {
				if literal, ok := item.(*ast.LiteralType); ok {
					results.Values = append(results.Values, literal.Token.Value())
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

func Walk(astNode ast.Node, query *Query, action func(node ast.Node, queryNode Node) error) error {
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
		// Assume a match if the for loop didn't return
		// Assume Keys will always be len > 0
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
