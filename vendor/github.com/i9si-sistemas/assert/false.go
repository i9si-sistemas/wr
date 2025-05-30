package assert

// False checks if the provided boolean value is false.
// If the value is true, it reports a fatal error using the provided fatalMessage.
//
//	result := someFunction()
//	assert.False(t, result, "Expected result to be false")
func False(t T, ok bool, args ...any) {
	tester := initTest(t)
	configureTest(tester, ok, false)
	if ok {
		tester.Fatal(args...)
	}
}
