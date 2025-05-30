package assert

import (
	"fmt"
	"reflect"
)

// Equal checks if the provided result is equal to the expected value.
// If the result is not equal to the expected value, it reports a fatal error
// using the provided fatalMessage.
//
//	var result int = someFunction()
//	expected := 42
//	assert.Equal(t, result, expected, "Expected result to be 42")
func Equal(t T, result, expected any, args ...any) {
	tester := initTest(t)
	configureTest(tester, result, expected)
	if !equal(result, expected) {
		tester.Fatal(args...)
	}
}

// StrictEqual checks if the provided result is equal to the expected value.
// If the result is not equal to the expected value, it reports a fatal error
// using the provided fatalMessage.
//
//	type value int
//	assert.StrictEqual(t, value(10), 10, "not equal")// not equal
//	assert.StrictEqual(t, value(20), value(10), "not equal")// not equal
//	assert.StrictEqual(t, value(20), value(20), "not equal")// pass
func StrictEqual(t T, result, expected any, args ...any) {
	tester := initTest(t)
	configureTest(tester, result, expected)
	if !strictEqual(result, expected) {
		tester.Fatal(args...)
	}
}

func strictEqual(result, expected any) bool {
	return reflect.DeepEqual(result, expected)
}

func equal(result, expected any) bool {
	return fmt.Sprint(result) == fmt.Sprint(expected)
}
