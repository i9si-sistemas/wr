package assert

// True checks if the provided boolean value is true.
// If the value is false, it reports a fatal error using the provided fatalMessage.
//
//	result := someFunction()
//	assert.True(t, result, "Expected result to be true")
func True(t T, ok bool, args ...any) {
	tester := initTest(t)
	configureTest(tester, ok, true)
	if !ok {
		tester.Fatal(args...)
	}
}
