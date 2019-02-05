package hclq

import (
	"fmt"

	"github.com/hashicorp/hcl/hcl/ast"
	jsonParser "github.com/hashicorp/hcl/json/parser"
	"github.com/valyala/fastjson"
)

// HclFromJSON converts a JSON string into an HCL ast.Node.
func HclFromJSON(str string) (ast.Node, error) {
	// hcl's built-in JSON parser doesn't handle invalid input very well,
	// so use fastjson instead.
	err := fastjson.Validate(str)
	if err != nil {
		return nil, fmt.Errorf("new value is not valid JSON: %s", err.Error())
	}

	str = fmt.Sprintf(`{"root":%s}`, str) // parser requires `root` key
	data := []byte(str)
	hcl, err := jsonParser.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("new value is not valid JSON: %s", err.Error())
	}
	node := hcl.Node.(*ast.ObjectList).Items[0].Val
	return node, nil
}

// HclListFromJSON the same as HclFromJSON but converts the result to a list.
func HclListFromJSON(str string) (*ast.ListType, error) {
	node, err := HclFromJSON(str)
	if err != nil {
		return nil, err
	}
	listNode, ok := node.(*ast.ListType)
	if !ok {
		return nil, fmt.Errorf("new value is not a list")
	}
	return listNode, nil
}

// HclLiteralFromJSON the same as HclFromJSON but converts the result to a literal value.
func HclLiteralFromJSON(str string) (*ast.LiteralType, error) {
	node, err := HclFromJSON(str)
	if err != nil {
		return nil, nil
	}
	literalNode, ok := node.(*ast.LiteralType)
	if !ok {
		return nil, fmt.Errorf("new value is not a literal")
	}
	return literalNode, nil
}
