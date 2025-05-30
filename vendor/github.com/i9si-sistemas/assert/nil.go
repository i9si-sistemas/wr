package assert

// Nil checks if the provided value is nil.
// If the value is not nil, it reports a fatal error using the provided fatalMessage.
//
//	myVar := new(string)
//	assert.Nil(t, myVar, "myVar should be nil")
func Nil(t T, v any, args ...any) {
	tester := initTest(t)
	configureTest(tester, v, "nil")
	if v != nil {
		tester.Fatal(args...)
	}
}
