package hclq

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
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

// Query returns results that match a given HCL2 expression
func (f *Doc) Query(expression string) hcl.Diagnostics {
	return nil
}

// FromFile creates a queryable HCL document from a filename.
func FromFile(path string) (*Doc, hcl.Diagnostics) {
	return &Doc{File: nil}, nil
}

// FromString creates a queryable HCL document from string contents.
func FromString(contents string) (result *Doc, errors hcl.Diagnostics) {
	file, diags := hclsyntax.ParseConfig([]byte(contents), "nofile", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags
	}
	return &Doc{File: file}, nil
}
