package assert
// Empty checks if the provided value is equal to the zero value of its type.
// If the value is not equal to the zero value, it reports a fatal error
// using the provided fatalMessage.
//
//	var myVar string
//	assert.Empty(t, myVar, "myVar should be empty string")
func Empty(t T, value any, args ...any){
	Zero(t, value, args...)
}