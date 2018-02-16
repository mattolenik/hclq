package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
)

func getOutput(obj interface{}, raw bool) (string, error) {
	if raw {
		result := ""
		arr, isArray := obj.([]interface{})
		if isArray {
			for _, item := range arr {
				result += fmt.Sprintf("%+v", item) + " "
			}
			result = strings.TrimRight(result, " ")
			result += "\n"
			return result, nil
		}
		return fmt.Sprintf("%+v\n", obj), nil
	}
	jsonBody, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonBody), nil
}
