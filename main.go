package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/k0kubun/pp"
	"github.com/mattolenik/hclq/queryast"
)

func main() {
	err := mainError()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var config string = `
a = {
	b = 4
	c = {
		f = "xz"
	}
}
`

var configBytes []byte = []byte(config)

func mainError() error {

	q := ".a.b | c()"
	r, err := queryast.Parse("inline", []byte(q))
	pp.Println(r)
	pp.Println(err)
	c, diags := hclsyntax.ParseConfig([]byte(config), "config", hcl.Pos{})
	if diags != nil {
		return diags
	}
	_, err = traverse(c.Body)
	pp.Println(c.Body)
	pp.Println(err)
	return nil
}

func extractRange(r hcl.Range) string {
	sc := hcl.NewRangeScannerFragment(configBytes, "config", r.Start, bufio.ScanLines)
	result := ""
	for sc.Scan() {
		if sc.Range() == r {
			result += string(sc.Bytes())
		}
	}
	return result
}

func traverseQuery(queryNode interface{}, hclNode interface{}) (interface{}, error) {
	if queryNode == nil {
	}
	if hclNode == nil {
	}
	switch n := queryNode.(type) {
	case *queryast.Expr:
		r, err := traverseQuery(n.Node, hclNode)
		if err != nil {
			return nil, err
		}
		// traverseQuery(n.Next) somehow pipe input here
		return r, nil
	}
	return nil, nil
}

func matchPath(crumbs []*queryast.Crumb, crumbIndex int, hclNode interface{}) (interface{}, error) {
	if crumbIndex > len(crumbs) {
		return hclNode, nil
	}
	crumb := crumbs[crumbIndex]
	switch crumb.Key.Selector.(type) {
	case *queryast.EmptySelector:
	case *queryast.IndexSelector:
	case *queryast.SplatSelector:
		//case *queryast.FunctionCall:
	default:
		fmt.Println("unexpected selector type in matchPath")
	}
	switch n := hclNode.(type) {
	case *hclsyntax.Body:
		attrNode, err := matchPath(crumbs, crumbIndex+1, n.Attributes)
		if err != nil {
			return nil, err
		}
		if attrNode != nil {
			return attrNode, nil
		}
		blockNode, err := matchPath(crumbs, crumbIndex+1, n.Blocks)
		if err != nil {
			return nil, err
		}
		if blockNode != nil {
			return blockNode, nil
		}
		return n, nil
	case *hclsyntax.Block:
		if n.Type != crumb.Key.Ident {
			return nil, nil
		}
		i, label := 0, ""
		for i, label = range n.Labels {
			if i+crumbIndex >= len(crumbs) {
				return nil, nil
			}
			if label != crumbs[i+crumbIndex].Key.Ident {
				return nil, nil
			}
		}
		return matchPath(crumbs, i+crumbIndex, n.Body)
	case hclsyntax.Blocks:
		for _, block := range n {
			blockNode, err := matchPath(crumbs, crumbIndex, block)
			if err != nil {
				return nil, err
			}
			if blockNode != nil {
				return blockNode, nil
			}
		}
		return nil, nil
	case *hclsyntax.Attribute:
		//n.Expr.Value()
		//n.Expr
		if n.Name == crumb.Key.Ident {
			return matchPath(crumbs, crumbIndex+1, n)
		}
		return nil, nil
	case hclsyntax.Attributes:
		if val, ok := n[crumb.Key.Ident]; ok {
			return matchPath(crumbs, crumbIndex+1, val)
		}
		return nil, nil
	default:
		fmt.Println("unhandled type in matchPath")
		return nil, nil
	}
}

func traverse(node interface{}) (interface{}, error) {
	switch v := node.(type) {
	case *hclsyntax.Body:
		for _, attr := range v.Attributes {
			_, err := traverse(attr)
			if err != nil {
				return nil, err
			}
		}
		for _, block := range v.Blocks {
			traverse(block)
		}
		return nil, nil
	case *hclsyntax.Attribute:
		return traverse(v.Expr)
	case *hclsyntax.Block:
		return traverse(v.Body)
	case *hclsyntax.ObjectConsExpr:
	case *hclsyntax.ObjectConsKeyExpr:
	case *hclsyntax.ScopeTraversalExpr:
	case *hclsyntax.TupleConsExpr:
	case hclsyntax.Expression:
		r := v.Range()
		fmt.Println(extractRange(r))
	default:
		fmt.Println(reflect.TypeOf(v))
	}
	return nil, nil
}
