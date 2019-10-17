package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/k0kubun/pp"
)

func main() {
	err := mainError()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mainError() error {
	config := `
type1 label1 label2 label3 {
a = upper(upper("abc"))
b lb1 {
d = "abc"
}
}
	`
	//q := os.Args[1]
	//_, err := queryast.Parse("inline", []byte(q))
	//pp.Println(r)
	//pp.Println(err)
	c, diags := hclsyntax.ParseConfig([]byte(config), "config", hcl.Pos{})
	if diags != nil {
		return diags
	}
	r, err := traverse(c.Body)
	pp.Println(r)
	pp.Println(err)
	return nil
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
		fmt.Println("a:" + v.Name)
		return traverse(v.Expr)
	case *hclsyntax.Block:
		fmt.Println("b:" + v.Type)
		return traverse(v.Body)
	case *hclsyntax.TemplateExpr:
	case *hclsyntax.FunctionCallExpr:
	default:
		fmt.Println(reflect.TypeOf(v))
	}
	return nil, nil
}
