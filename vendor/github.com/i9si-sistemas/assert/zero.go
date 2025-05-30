package assert

import (
	"reflect"
)

// Zero checks if the provided value is equal to the zero value of its type.
// If the value is not equal to the zero value, it reports a fatal error
// using the provided fatalMessage.
//
//	var myVar string
//	assert.Zero(t, myVar, "myVar should be empty string")
func Zero(t T, value any, args ...any) {
	tester := initTest(t)
	zeroValue, ok := isZeroValue(value)
	configureTest(tester, value, zeroValue)
	if !ok {
		tester.Fatal(args...)
	}
}

func isZeroValue(value any) (v reflect.Value, ok bool) {
	v = reflect.Zero(reflect.TypeOf(value))
	ok = equal(value, v.Interface())
	return
}
