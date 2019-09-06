package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/mattolenik/hclq/hclq"
	"github.com/zclconf/go-cty/cty"
)

var version string

var hcl = `
some "block" "here" {
    foo="abc"
}
`

func main() {
	doc, errs := hclq.FromString(hcl)
	if errs != nil {
		panic(errs)
	}
	spec := &hcldec.BlockObjectSpec{
		TypeName:   "some",
		LabelNames: []string{"block", "here"},
		Nested: &hcldec.ObjectSpec{
			"foo": &hcldec.AttrSpec{Name: "foo", Type: cty.String},
		},
	}
	res, err := doc.Query(spec)
	if err != nil {
		panic(err)
	}
	spew.Dump(res)
}
