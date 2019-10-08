package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/repr"
	"github.com/mattolenik/hclq/queryast"
)

func main() {
	err := mainError()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mainError() error {
	q := ".abc.Ǔdef.xyz"
	r, err := queryast.Parse("inline", []byte(q))
	repr.Println(r)
	repr.Println(err)
	return nil
}
