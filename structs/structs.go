package structs

// AsSlice casts v to an array of interface{} from just an interface{},
// returning nil if v is nil.
func AsSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}
