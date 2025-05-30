package assert

// NoError checks if the provided error is nil.
// If the error is not nil, it reports a fatal error using the provided fatalMessage.
//
//	err := someFunction()
//	fatalMessage := fmt.Sprintf("Expected no error, but got: %v", err)
//	assert.NoError(t, err, fatalMessage)
func NoError(t T, err error, args ...any) {
	tester := initTest(t)
	configureTest(tester, err, "no error")
	Nil(t, err, args...)
}
