package hclq

import (
	"fmt"
	"strconv"

    "github.com/hashicorp/hcl/hcl/ast"
)

func (doc *HclDocument) GetRaw(q string) ([]ast.Node, error) {
	resultPairs, err := doc.Query(q)
	if err != nil {
		return nil, err
	}
	results := []ast.Node{}
	for _, pair := range resultPairs {
		results = append(results, pair.Node)
	}
	return results, nil
}

// Get performs a query and returns a deserialized value. The query string is the same format as the command line.
func (doc *HclDocument) Get(q string) (interface{}, error) {
	resultPairs, err := doc.Query(q)
	if err != nil {
		return nil, err
	}
	results := make([]interface{}, len(resultPairs))
	for _, pair := range resultPairs {
		results = append(results, pair.Value)
	}
	return results, nil
}

// GetKeys peforms a query and returns just the key or keys, no values.
func (doc *HclDocument) GetKeys(q string) ([]string, error) {
	return doc.QueryKeys(q)
}

// GetAsInt performs Get but converts the result to a string
func (doc *HclDocument) GetAsInt(q string) (int, error) {
	result, err := doc.Get(q)
	if err != nil {
		return 0, err
	}
	num, ok := result.(int)
	if ok {
		return num, nil
	}
	str, ok := result.(string)
	if ok {
		num, err := strconv.Atoi(str)
		if err == nil {
			return num, nil
		}
	}
	return 0, fmt.Errorf("could not find int at '%s' nor a string convertable to an int", q)
}

// GetAsString performs Get but converts the result to a string
func (doc *HclDocument) GetAsString(q string) (string, error) {
	result, err := doc.Get(q)
	if err != nil {
		return "", err
	}
	str, ok := result.(string)
	if ok {
		return str, nil
	}
	num, ok := result.(int)
	if ok {
		return strconv.Itoa(num), nil
	}
	return fmt.Sprintf("%v", result), nil
}

// GetAsList does the same as Get but converts it to a list for you (with type check)
func (doc *HclDocument) GetAsList(q string) ([]interface{}, error) {
	result, err := doc.Get(q)
	if err != nil {
		return nil, err
	}
	arr, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("query does not return a list, cannot be used with GetList")
	}
	return arr, nil
}

// GetAsStringList does the same as GetAsList but converts everything to a string for you.
func (doc *HclDocument) GetAsStringList(q string) ([]string, error) {
	list, err := doc.GetAsList(q)
	if err != nil {
		return nil, err
	}
	results := make([]string, len(list))
	for i, item := range list {
		str, ok := item.(string)
		if ok {
			results[i] = str
			continue
		}
		num, ok := item.(int)
		if ok {
			results[i] = strconv.Itoa(num)
			continue
		}
		// Fall back to general Go print formatting
		results[i] = fmt.Sprintf("%v", item)
	}
	return results, nil
}

// GetAsIntList does the same as GetAsList but with all values converted to ints.
// Returns an error if a value is found that is not an int and couldn't be parsed into one.
func (doc *HclDocument) GetAsIntList(q string) ([]int, error) {
	list, err := doc.GetAsList(q)
	if err != nil {
		return nil, err
	}
	results := make([]int, len(list))
	for i, item := range list {
		num, ok := item.(int)
		if ok {
			results[i] = num
			continue
		}
		str, ok := item.(string)
		if ok {
			num, err := strconv.Atoi(str)
			if err != nil {
				return nil, fmt.Errorf("failed to parse '%s' into an integer", str)
			}
			results[i] = num
			continue
		}
		return nil, fmt.Errorf("value '%v' is not an integer and could not be parsed into one", item)
	}
	return results, nil
}
