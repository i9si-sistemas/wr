package assert 

// NotEmpty checks if the provided value is not equal to the zero value of its type.
// If the value is equal to the zero value, it reports a fatal error
// using the provided fatalMessage.
//
//	var myVar string
//	assert.NotEmpty(t, myVar, "myVar should not be empty string")
func NotEmpty(t T, value any, args ...any){
	tester := initTest(t)
	zeroValue, ok := isZeroValue(value)
	configureTest(tester, zeroValue, "should not be zero value")
	if ok {
		t.Fatal(args...)
	}
}