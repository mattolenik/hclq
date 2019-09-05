package hclq

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// Doc is a queryable HCL file
type Doc struct {
	File *hcl.File
}

// Results contains HCLQ query results
type Results struct {
	Value cty.Value
}

// Query returns results that match a given query string
func (f *Doc) Query(query string) (*Results, error) {
	traverse := hcl.TraverseAttr{Name: "foo", SrcRange: nil}
	return nil, nil
}

// FromFile creates a queryable HCL document from a filename.
func FromFile(path string) (*Doc, hcl.Diagnostics) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(path)
	if diags != nil {
		return nil, diags
	}
	return &Doc{File: file}, nil
}

// FromString creates a queryable HCL document from string contents.
func FromString(contents string) (result *Doc, errors hcl.Diagnostics) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(contents), "nofile")
	if diags != nil {
		return nil, diags
	}
	return &Doc{File: file}, nil
}
