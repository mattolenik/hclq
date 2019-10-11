package structs

import (
	"fmt"
	"reflect"
)

// AsSlice casts v to an array of interface{} from just an interface{},
// returning nil if v is nil.
func AsSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

// GetOfType takes a slice and returns the first element matching typ, or nil if there is none.
func GetOfType(slice interface{}, typ reflect.Type) interface{} {
	if slice == nil {
		return nil
	}
	sl, ok := slice.([]interface{})
	if !ok {
		panic("expected slice type")
	}
	for _, v := range sl {
		if reflect.TypeOf(v) == typ {
			return v
		}
	}
	return nil
}

// SingleOrNil converts a single-item slice into a single value. If v is already a single value,
// it is simply returned as-is. If the slice v has more than one element, nil and an error are returned.
// If passed nil, nil will be returned.
func SingleOrNil(v interface{}) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	if slice, ok := v.([]interface{}); ok {
		if len(slice) == 1 {
			return slice[0], nil
		} else {
			return nil, fmt.Errorf("argument v must be a slice of length 1 or not a slice at all")
		}
	}
	return v, nil
}
