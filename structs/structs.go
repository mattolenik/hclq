package structs

import (
	"fmt"
	"reflect"
	"strconv"
)

var X = 1

// Fill takes a pointer to a struct and populates fields tagged with 'index:n' where n is an integer >= 0.
// That index corresponds to the index of the values slice that is passed into this function. A field with
// the tag `index:"1"` will be populated with values[1]. This allows fields to be rearranged in code without
// disturbing the fill order.
//
// type Fillable struct {
//     A int    `index:"0"`
//     B string `index:"1"`
// }
//
// f := &Fillable{}
// structs.Fill(f, 123, "xyz")
// // f now equals {A: 123, B: "xyz"}
func Fill(record interface{}, values []interface{}) error {
	tagName := "index"
	recordValue := reflect.ValueOf(record).Elem()
	recordType := reflect.TypeOf(record).Elem()
	numValues := len(values)
	numFields := recordValue.NumField()
	if numValues > numFields {
		return fmt.Errorf("tried to fill %d values, but struct only has %d fields", numValues, numFields)
	}
	fieldIndexToOrder := map[int]int{}
	orderToFieldIndex := map[int]int{}
	for i := 0; i < numFields; i++ {
		field := recordValue.Field(i)
		t, ok := recordType.Field(i).Tag.Lookup(tagName)
		// Field does not have order tag, skip it
		if !ok {
			continue
		}
		order, err := strconv.Atoi(t)
		if err != nil || order < 0 {
			return fmt.Errorf("%s tag must be a valid integer greater than or equal to 0", tagName)
		}
		if order >= numValues {
			return fmt.Errorf("%s %d on field named '%s' is out of bounds, only %d values were provided", tagName, order, recordType.Field(i).Name, numValues)
		}
		if fi, ok := orderToFieldIndex[order]; ok {
			return fmt.Errorf("%s %d already used on field number %d", tagName, order, fi)
		}
		if !field.IsValid() {
			return fmt.Errorf("field at index %d not a valid value", i)
		}
		fieldIndexToOrder[i] = order
		orderToFieldIndex[order] = i
	}
	for index, order := range fieldIndexToOrder {
		field := recordValue.Field(index)
		field.Set(reflect.ValueOf(values[order]))
	}
	return nil
}
