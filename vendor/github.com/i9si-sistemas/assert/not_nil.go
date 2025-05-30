package assert

// NotNil checks if the provided value is not nil.
// If the value is nil, it reports a fatal error using the provided fatalMessage.
//
//	var myVar *string
//	assert.NotNil(t, myVar, "myVar should not be nil")
func NotNil(t T, v any, args ...any) {
	tester := initTest(t)
	configureTest(tester, v, "not nil")
	if v == nil {
		tester.Fatal(args...)
	}
}
