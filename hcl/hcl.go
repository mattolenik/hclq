package hcl

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/mattolenik/hclq/query"
	//"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/parser"
)

// Result represents a query result
type Result struct {
	Value interface{}
	Node  ast.Node
}

// Get performs a query and returns a deserialized value
func Get(reader io.Reader, q string) (interface{}, error) {
	qry, _ := query.Parse(q)
	resultPairs, isList, _, err := Query(reader, qry)
	if err != nil {
		return nil, err
	}
	results := []interface{}{}
	for _, pair := range resultPairs {
		results = append(results, pair.Value)
	}
	// The return type can be a list if: the queried object IS a list, or if the query matched multiple single items
	// So, return now if it's not a list and there is only one query result
	if !isList && len(results) == 1 {
		return results[0], nil
	}
	return results, nil
}

// GetAsInt performs Get but converts the result to a string
func GetAsInt(reader io.Reader, q string) (int, error) {
	result, err := Get(reader, q)
	if err != nil {
		return 0, err
	}
	num, ok := result.(int)
	if ok {
		return num, nil
	}
	str, ok := result.(string)
	if ok {
		num, err := strconv.Atoi(str)
		if err == nil {
			return num, nil
		}
	}
	return 0, fmt.Errorf("Could not find int at '%s' nor a string convertable to an int", q)
}

// GetAsString performs Get but converts the result to a string
func GetAsString(reader io.Reader, q string) (string, error) {
	result, err := Get(reader, q)
	if err != nil {
		return "", err
	}
	str, ok := result.(string)
	if ok {
		return str, nil
	}
	num, ok := result.(int)
	if ok {
		return strconv.Itoa(num), nil
	}
	return fmt.Sprintf("%v", result), nil
}

// GetAsList does the same as Get but converts it to a list for you (with type check)
func GetAsList(reader io.Reader, q string) ([]interface{}, error) {
	result, err := Get(reader, q)
	if err != nil {
		return nil, err
	}
	arr, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Query does not return a list, cannot be used with GetList")
	}
	return arr, nil
}

// GetAsStringList does the same as GetAsList but converts everything to a string for you.
func GetAsStringList(reader io.Reader, q string) ([]string, error) {
	list, err := GetAsList(reader, q)
	if err != nil {
		return nil, err
	}
	results := make([]string, len(list))
	for _, item := range list {
		str, ok := item.(string)
		if ok {
			results = append(results, str)
			continue
		}
		num, ok := item.(int)
		if ok {
			results = append(results, strconv.Itoa(num))
			continue
		}
		// Fall back to general Go print formatting
		results = append(results, fmt.Sprintf("%v", item))
	}
	return results, nil
}

// GetAsIntList does the same as GetAsList but with all values converted to ints.
// Returns an error if a value is found that is not an int and couldn't be parsed into one.
func GetAsIntList(reader io.Reader, q string) ([]int, error) {
	list, err := GetAsList(reader, q)
	if err != nil {
		return nil, err
	}
	results := make([]int, len(list))
	for _, item := range list {
		num, ok := item.(int)
		if ok {
			results = append(results, num)
			continue
		}
		str, ok := item.(string)
		if ok {
			num, err := strconv.Atoi(str)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse '%s' into an integer", str)
			}
			results = append(results, num)
			continue
		}
		return nil, fmt.Errorf("value '%v' is not an integer and could not be parsed into one", item)
	}
	return results, nil
}

// Query performs a generic query and returns matching results
func Query(reader io.Reader, qry *query.Breadcrumbs) (results []Result, isList bool, node *ast.File, err error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	node, err = parser.Parse(bytes)
	if err != nil {
		return
	}
	err = walk(node.Node, qry, 0, func(n ast.Node, crumb query.Crumb) error {
		switch node := n.(type) {
		case *ast.LiteralType:
			results = append(results, Result{node.Token.Value(), node})

		case *ast.ListType:
			listNode, ok := crumb.(query.IndexedCrumb)
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

func walk(astNode ast.Node, query *query.Breadcrumbs, qIdx int, action func(node ast.Node, crumb query.Crumb) error) error {
	switch node := astNode.(type) {
	case *ast.ObjectList:
		for _, obj := range node.Items {
			err := walk(obj, query, qIdx, action)
			if err != nil {
				return err
			}
		}
		return nil

	case *ast.ObjectItem:
		for _, key := range node.Keys {
			part := query.Parts[qIdx]
			isMatch, err := part.IsMatch(strings.Trim(key.Token.Text, `"`), node.Val)
			if err != nil {
				return err
			}
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
		return walk(node.Val, query, qIdx, action)

	case *ast.ListType:
		return action(node, query.Parts[qIdx])

	case *ast.LiteralType:
		return action(node, query.Parts[qIdx])

	case *ast.ObjectType:
		return walk(node.List, query, qIdx, action)

	default:
		return errors.New("unhandled case")
	}
}
