package main

import (
	"bufio"
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

var config string = `
type1 label1 label2 label3 {
a = upper(upper("abc"))
b lb1 {
d = "abc"
}
}
`

var configBytes []byte = []byte(config)

func mainError() error {

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
	case hclsyntax.Expression:
		r := v.Range()
		fmt.Println(extractRange(r))
	default:
		fmt.Println(reflect.TypeOf(v))
	}
	return nil, nil
}
