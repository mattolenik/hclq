package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/mattolenik/hclq/config"
	"io"
	"os"
	"strings"

    "github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl"
    "github.com/hashicorp/hcl/hcl/printer"
	//"github.com/davecgh/go-spew/spew"
)

type PrintStyle int
type MergeMethod int

const (
	HCL PrintStyle = iota
	JSON
	Raw
)

const (
	Combined = iota
	Separate
)

func PrintHCL(writer io.Writer, style PrintStyle, method MergeMethod, nodes ...ast.Node) error {
	var node ast.Node
	if len(nodes) == 0 {
		return fmt.Errorf("must provide at least one ast.Node")
	} else if len(nodes) > 1 {
		node = &ast.ListType { List: nodes, }
	} else {
		node = nodes[0]
	}
	switch style {
	case JSON:
		var decoded interface{}
		if len(nodes) > 1 {
			decoded = []map[string]interface{}{}
		} else {
			decoded = map[string]interface{}{}
		}
		err := hcl.DecodeObject(&decoded, node)
		if err != nil {
			return err
		}
		json, err := json.Marshal(&decoded)
		if err != nil {
			return err
		}
		_, err = writer.Write(json)
		writer.Write([]byte("\n"))
		if err != nil {
			return err
		}
	case HCL:
		for _, node := range nodes {
			err := printer.Fprint(writer, node)
			writer.Write([]byte("\n"))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getOutput(obj interface{}, raw bool) (string, error) {
	if raw {
		result := ""
		arr, isArray := obj.([]interface{})
		if isArray {
			for _, item := range arr {
				// Simple output, uses built-in %v, most useful for simple types.
				result += strings.Trim(fmt.Sprintf("%v", item), `"`) + " "
			}
			result = strings.TrimRight(result, " ")
			return result, nil
		}
		return strings.Trim(fmt.Sprintf("%v", obj), `"`), nil
	}
	jsonBody, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonBody), nil
}

// getInputReader provides an io.Reader for reading from either a file
// or stdin, depending on whether or not an input file was specified.
func getInputReader() (io.Reader, error) {
	if val := config.InputFile; val != "" {
		reader, err := os.Open(val)
		if err != nil {
			return nil, err
		}
		return reader, nil
	}
	return os.Stdin, nil
}
