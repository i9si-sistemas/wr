package assert

// Error checks if the provided error is not nil.
// If the error is nil, it reports a fatal error using the provided fatalMessage.
//
//	err := someFunction()
//	assert.Error(t, err, "Expected an error, but got nil")
func Error(t T, err error, args ...any) {
	tester := initTest(t)
	configureTest(tester, err, "error")
	NotNil(t, err, args...)
}
