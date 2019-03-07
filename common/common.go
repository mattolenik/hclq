package common

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// TODO: refactor, move these into more appropriate package

func GetOutput(obj interface{}, raw bool) (string, error) {
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

func GetInputReader(inFile *string) (io.Reader, error) {
	if *inFile != "" {
		reader, err := os.Open(*inFile)
		if err != nil {
			return nil, err
		}
		return reader, nil
	}
	return os.Stdin, nil
}

func GetOutputWriter(inFile, outFile *string, inPlace *bool) (io.WriteCloser, error) {
	if *outFile != "" {
		if *inPlace {
			return nil, fmt.Errorf("cannot use --in-place with --out")
		}
		return os.Create(*outFile)
	}
	if *inPlace {
		if *inFile != "" {
			return nil, fmt.Errorf("must specify input file when using --in-place")
		}
		return os.Create(*inFile)
	}
	return os.Stdout, nil
}
