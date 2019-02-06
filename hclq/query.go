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
	Key   string
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

// FromFile creates a new document from a file.
func FromFile(filename string) (*HclDocument, error) {
	bytes, err := ioutil.ReadFile(filename)
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

// QueryKeys performs a query and returns matching key or keys, no values.
func (doc *HclDocument) QueryKeys(queryString string) (keys []string, err error) {
	keys = []string{}
	results, err := doc.Query(queryString)
	if err != nil {
		return nil, err
	}
	for _, r := range results {
		if len(r.Key) > 0 {
			keys = append(keys, r.Key)
		}
	}
	return keys, nil
}

// Query performs a generic query and returns matching results
func (doc *HclDocument) Query(queryString string) (results []Result, err error) {
	qry, err := query.ParseBreadcrumbs(queryString)
	if err != nil {
		return nil, err
	}
	err = walk(doc.FileNode.Node, qry, []string{}, 0, func(astNode ast.Node, key string, crumb query.Crumb) error {
		switch node := astNode.(type) {
		case *ast.LiteralType:
			results = append(results, Result{key, node.Token.Value(), node})

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
				results = append(results, Result{key, val.Token.Value(), node})
				return nil
			}
			// Otherwise query is for all items. Add them to results as a new list.
			values := []interface{}{}
			for _, item := range node.List {
				if literal, ok := item.(*ast.LiteralType); ok {
					values = append(values, literal.Token.Value())
				}
			}
			results = append(results, Result{key, values, node})

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

func walk(
	astNode ast.Node,
	query *query.Breadcrumbs,
	keyTrail []string,
	crumbIndex int,
	action func(node ast.Node, key string, crumb query.Crumb) error) error {

	switch node := astNode.(type) {
	case *ast.ObjectList:
		for _, obj := range node.Items {
			err := walk(obj, query, keyTrail, crumbIndex, action)
			if err != nil {
				return err
			}
		}
		return nil

	case *ast.ObjectItem:
		for _, key := range node.Keys {
			part := query.Parts[crumbIndex]
			token := strings.Trim(key.Token.Text, `"`)
			keyTrail = append(keyTrail, token)
			isMatch, err := part.IsMatch(token, node.Val)
			if err != nil {
				return err
			}
			if !isMatch {
				return nil
			}
			// End of query, this means a match
			if crumbIndex+1 >= query.Length {
				break
			}
			crumbIndex++
		}
		// Assume a match if return didn't happen in the for loop.
		// Assume Keys will always be len > 0 (it wouldn't be valid HCL otherwise)
		return walk(node.Val, query, keyTrail, crumbIndex, action)

	case *ast.ListType:
		return action(node, strings.Join(keyTrail, "."), query.Parts[crumbIndex])

	case *ast.LiteralType:
		return action(node, strings.Join(keyTrail, "."), query.Parts[crumbIndex])

	case *ast.ObjectType:
		return walk(node.List, query, keyTrail, crumbIndex, action)

	default:
		return errors.New("unhandled case")
	}
}
