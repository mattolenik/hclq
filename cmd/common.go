package cmd

import (
	"encoding/json"
	"fmt"
)

func getOutput(obj interface{}, raw bool) (string, error) {
	if raw {
		return fmt.Sprintf("%+v", obj), nil
	}
	jsonBody, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("failure while trying to serialize output to JSON")
	}
	return string(jsonBody), nil
}
