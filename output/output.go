package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	hclPrinter "github.com/hashicorp/hcl/hcl/printer"
	//"github.com/davecgh/go-spew/spew"
)

type PrintStyle int

const (
	HCL PrintStyle = iota
	JSON
	Raw
)

func PrintStyleFromString(s string) (PrintStyle, error) {
	switch s {
	case "hcl":
		return HCL, nil
	case "json":
		return JSON, nil
	case "raw":
		return Raw, nil
	default:
		return -1, fmt.Errorf("invalid output format, must be one of: hcl, json, raw")
	}
}

func decode(node ast.Node) (interface{}, error) {
	var decoded interface{}
	switch n := node.(type) {
	case *ast.ObjectItem:
		// If an ObjectItem's Keys is length 1, it must be an assignment.
		if len(n.Keys) == 1 {
			err := hcl.DecodeObject(&decoded, n.Val)
			if err != nil {
				return nil, err
			}
		} else {
			decoded = map[string]interface{}{}
			err := hcl.DecodeObject(&decoded, n)
			if err != nil {
				return nil, err
			}
		}
	default:
		decoded := map[string]interface{}{}
		err := hcl.DecodeObject(&decoded, node)
		if err != nil {
			return nil, err
		}
	}
	return decoded, nil
}

func Print(writer io.Writer, style PrintStyle, nodes ...ast.Node) error {
	for _, node := range nodes {
		switch style {
		case JSON:
			decoded, err := decode(node)
			if err != nil {
				return err
			}
			printJSON(writer, decoded)
			fmt.Fprintf(writer, "\n")

		case Raw:
			decoded, err := decode(node)
			if err != nil {
				return err
			}
			fmt.Fprintf(writer, "%v\n", decoded)

		case HCL:
			hclPrinter.Fprint(writer, node)
			fmt.Fprintf(writer, "\n\n")
		}
	}
	return nil
}

func printJSON(writer io.Writer, obj interface{}) error {
	json, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = writer.Write(json)
	if err != nil {
		return err
	}
	return nil
}
