package utils

import (
	"errors"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
	//JsonParser "github.com/hashicorp/hcl/json/parser"
	"github.com/mattolenik/hclq/query"
	"strings"
)

type WalkAction func(node ast.Node) (stop bool, err error)

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

func Walk(node ast.Node, query []query.Node, queryIdx int, action WalkAction) (bool, error) {
	if objList, ok := node.(*ast.ObjectList); ok {
		for _, obj := range objList.Items {
			stop, err := Walk(obj, query, queryIdx, action)
			if err != nil {
				return stop, err
			}
			if stop {
				return stop, nil
			}
		}
		return false, nil
	}
	if objItem, ok := node.(*ast.ObjectItem); ok {
		queryLen := len(query)
		for _, key := range objItem.Keys {
			if queryIdx >= queryLen {
				return false, nil
			}
			value := strings.Trim(key.Token.Text, "\"")
			if !query[queryIdx].IsMatch(value) {
				return false, nil
			}
			queryIdx++
		}
		// Assume a match if the for loop didn't return
		// Assume Keys will always be len > 0
		return Walk(objItem.Val, query, queryIdx, action)
	}
	switch _ := node.(type) {
	case *ast.ListType:
		return action(node)
	case *ast.LiteralType:
		return action(node)
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
	if objType, ok := node.(*ast.ObjectType); ok {
		return Walk(objType.List, query, queryIdx, action)
	}
	return false, errors.New("unhandled case")
}