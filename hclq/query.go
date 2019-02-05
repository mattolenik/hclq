package hclq

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/mattolenik/hclq/query"
)

// Result represents a query result
type Result struct {
	Value interface{}
	Node  ast.Node
}

// HclDocument represents an HCL document in memory.
type HclDocument struct {
	FileNode *ast.File
}

// FromReader creates a new document from an io.Reader
func FromReader(reader io.Reader) (*HclDocument, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	doc := &HclDocument{}
	doc.FileNode, err = parser.Parse(bytes)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// Print writes the HCL document out to the given io.Writer.
func (doc *HclDocument) Print(writer io.Writer) error {
	return printer.Fprint(writer, doc.FileNode)
}

// Query performs a generic query and returns matching results
func (doc *HclDocument) Query(queryString string) (results []Result, err error) {
	qry, err := query.ParseBreadcrumbs(queryString)
	if err != nil {
		return nil, err
	}
	err = walk(doc.FileNode.Node, qry, 0, func(astNode ast.Node, crumb query.Crumb) error {
		switch node := astNode.(type) {
		case *ast.LiteralType:
			results = append(results, Result{node.Token.Value(), node})

		case *ast.ListType:
			listNode, ok := crumb.(query.IndexedCrumb)
			if !ok {
				return fmt.Errorf("unexpected query type")
			}
			// In this case, the query is for a specific index. Add it to results as a single item.
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
			// Otherwise query is for all items. Add them to results as a new list.
			// TODO: should be done recursively to handle sub-lists
			values := []interface{}{}
			for _, item := range node.List {
				if literal, ok := item.(*ast.LiteralType); ok {
					values = append(values, literal.Token.Value())
				}
			}
			results = append(results, Result{Value: values, Node: node})
			return nil
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
