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

func printJSON(writer io.Writer, obj interface{}) error {
	json, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = writer.Write(json)
	writer.Write([]byte("\n"))
	if err != nil {
		return err
	}
	return nil
}

func Print(writer io.Writer, style PrintStyle, method MergeMethod, nodes ...ast.Node) error {
	for _, node := range nodes {
		switch style {
		case JSON:
			var decoded interface{}
			switch n := node.(type) {
			case *ast.ObjectItem:
				// If an ObjectItem's Keys is length 1, it must be an assignment.
				if len(n.Keys) == 1 {
					err := hcl.DecodeObject(&decoded, n.Val)
					if err != nil {
						return err
					}
				} else {
					decoded = map[string]interface{}{}
					err := hcl.DecodeObject(&decoded, n)
					if err != nil {
						return err
					}
				}
			default:
				decoded := map[string]interface{}{}
				err := hcl.DecodeObject(&decoded, node)
				if err != nil {
					return err
				}
			}
			printJSON(writer, decoded)

		case HCL:
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
