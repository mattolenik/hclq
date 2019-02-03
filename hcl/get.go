package hcl

import (
	"fmt"
	"github.com/mattolenik/hclq/query"
	"io"
	"strconv"
)

// Get performs a query and returns a deserialized value
func Get(reader io.Reader, q string) (interface{}, error) {
	qry, _ := query.Parse(q)
	resultPairs, isList, _, err := Query(reader, qry)
	if err != nil {
		return nil, err
	}
	results := []interface{}{}
	for _, pair := range resultPairs {
		results = append(results, pair.Value)
	}
	// The return type can be a list if: the queried object IS a list, or if the query matched multiple single items
	// So, return now if it's not a list and there is only one query result
	if !isList && len(results) == 1 {
		return results[0], nil
	}
	return results, nil
}

// GetAsInt performs Get but converts the result to a string
func GetAsInt(reader io.Reader, q string) (int, error) {
	result, err := Get(reader, q)
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
	return 0, fmt.Errorf("Could not find int at '%s' nor a string convertable to an int", q)
}

// GetAsString performs Get but converts the result to a string
func GetAsString(reader io.Reader, q string) (string, error) {
	result, err := Get(reader, q)
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
func GetAsList(reader io.Reader, q string) ([]interface{}, error) {
	result, err := Get(reader, q)
	if err != nil {
		return nil, err
	}
	arr, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Query does not return a list, cannot be used with GetList")
	}
	return arr, nil
}

// GetAsStringList does the same as GetAsList but converts everything to a string for you.
func GetAsStringList(reader io.Reader, q string) ([]string, error) {
	list, err := GetAsList(reader, q)
	if err != nil {
		return nil, err
	}
	results := make([]string, len(list))
	for _, item := range list {
		str, ok := item.(string)
		if ok {
			results = append(results, str)
			continue
		}
		num, ok := item.(int)
		if ok {
			results = append(results, strconv.Itoa(num))
			continue
		}
		// Fall back to general Go print formatting
		results = append(results, fmt.Sprintf("%v", item))
	}
	return results, nil
}

// GetAsIntList does the same as GetAsList but with all values converted to ints.
// Returns an error if a value is found that is not an int and couldn't be parsed into one.
func GetAsIntList(reader io.Reader, q string) ([]int, error) {
	list, err := GetAsList(reader, q)
	if err != nil {
		return nil, err
	}
	results := make([]int, len(list))
	for _, item := range list {
		num, ok := item.(int)
		if ok {
			results = append(results, num)
			continue
		}
		str, ok := item.(string)
		if ok {
			num, err := strconv.Atoi(str)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse '%s' into an integer", str)
			}
			results = append(results, num)
			continue
		}
		return nil, fmt.Errorf("value '%v' is not an integer and could not be parsed into one", item)
	}
	return results, nil
}
