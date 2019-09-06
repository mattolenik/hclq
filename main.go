package main

import (
	"fmt"

	"github.com/mattolenik/hclq/hclq"
)

var version string

func main() {
	doc, errs := hclq.FromString(`foo="123"`)
	if errs != nil {
		panic(errs)
	}
	res, err := doc.Query("foo")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
