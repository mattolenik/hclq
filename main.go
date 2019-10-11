package main

import (
	"fmt"
	"os"

	//"github.com/alecthomas/repr"
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

func mainError() error {
	q := os.Args[1]
	r, err := queryast.Parse("inline", []byte(q))
	pp.Println(r)
	pp.Println(err)
	return nil
}
