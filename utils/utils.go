package utils

import (
	"errors"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
	//JsonParser "github.com/hashicorp/hcl/json/parser"
	"github.com/mattolenik/hclq/query"
)

type WalkAction func(node ast.Node, queryNode query.Node) (stop bool, err error)

type ErrorJSON struct {
	Error string `json:"error"`
}

func ToGoType(node ast.Node) (interface{}, error) {
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
			nextItem, err := ToGoType(item)
			if err != nil {
				return nil, err
			}
			result = append(result, nextItem)
		}
		return result, nil
	} else if objectItem, ok := node.(*ast.ObjectItem); ok {
		result, err := ToGoType(objectItem.Val)
		return result, err
	}
	return "", errors.New("unhandled type conversion")
}

func Walk(astNode ast.Node, query []query.Node, queryIdx int, action WalkAction) (bool, error) {
	switch node := astNode.(type) {
	case *ast.ObjectList:
		for _, obj := range node.Items {
			stop, err := Walk(obj, query, queryIdx, action)
			if err != nil {
				return stop, err
			}
			if stop {
				return stop, nil
			}
		}
		return false, nil

	case *ast.ObjectItem:
		queryLen := len(query)
		for _, key := range node.Keys {
			if !query[queryIdx].IsMatch(key.Token.Text, node.Val) {
				return false, nil
			}
			if queryIdx + 1 >= queryLen {
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
		return false, errors.New("unhandled case")
	}
	//if list, ok := node.(*ast.ListType); ok {
		// HCL JSON parser needs a top level object
		//jsonValue := fmt.Sprintf(`{"root": %s}`, )
		//tree, err := JsonParser.Parse([]byte(jsonValue))
		//if err != nil {
		//	return err
		//}
		//list.List = tree.Node.(*ast.ObjectList).Items[0].Val.(*ast.ListType).List
		//return nil
	//}
}