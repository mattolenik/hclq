package main

import (
	"github.com/alecthomas/repr"
)

func main() {
	q := "abc"
	r, err := Parse("inline", []byte(q))
	repr.Println(r)
	repr.Println(err)
}
