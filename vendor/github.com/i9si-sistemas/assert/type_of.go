package assert

import "reflect"

// typeOf returns the name of the concrete type of the given value.
// For pointers, it dereferences and returns the base type name.
// Returns an empty string for nil values or unnamed types.
func typeOf(data any) string {
	if data == nil {
		return ""
	}

	t := reflect.TypeOf(data)
	if t == nil {
		return ""
	}

	// Dereference pointers
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Return the type name (empty for unnamed types)
	return t.Name()
}
