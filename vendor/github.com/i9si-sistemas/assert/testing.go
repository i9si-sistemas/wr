package assert

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/i9si-sistemas/assert/call"
)

// T is an interface that provides methods for managing test execution,
// reporting failures, and handling environment variables during tests.
type T interface {
	// Cleanup registers a function to be called when the test finishes.
	// This is useful for cleaning up resources that were allocated during the test.
	Cleanup(func())

	// Fail marks the function as having failed but continues execution.
	// This allows you to perform additional checks or logging even after a failure is detected.
	Fail()

	// Failed reports whether the function has failed.
	// It returns true if the test has been marked as failed by calling Fail().
	Failed() bool

	// Fatal logs the provided arguments as a fatal error and stops test execution immediately.
	// Unlike Fail(), this method will not allow the test to continue after it is called.
	Fatal(args ...any)

	// Helper marks the calling function as a helper function.
	Helper()

	// Setenv sets an environment variable for the duration of the test.
	// This can be useful for configuring the environment in which the test is run.
	Setenv(key, value string)

	// Skip marks the test as skipped and logs the provided arguments as the reason for skipping.
	// A skipped test will not be counted as a failure.
	Skip(args ...any)

	// Skipped reports whether the test has been skipped.
	// It returns true if the test was marked as skipped by calling Skip().
	Skipped() bool

	// TempDir returns the path to a temporary directory that can be used for storing temporary files
	// during the execution of the test. This directory is usually cleaned up after the test completes.
	TempDir() string
}

type test struct {
	ctx                              context.Context
	t                                T
	Caller                           string
	result, expected                 string
	originalResult, originalExpected any
	exit                             func(code int)
}

func newTest(ctx context.Context, t T) *test {
	return &test{ctx: ctx, t: t, exit: os.Exit}
}

func (t *test) setCaller(caller string) {
	t.Caller = caller
}

func (t *test) Context() context.Context {
	return t.ctx
}

func (t *test) Cleanup(f func()) {
	t.t.Cleanup(f)
}

func (t *test) Fail() {
	fmt.Print(t.failedMessage())
	t.t.Fail()
}

func (t *test) Failed() bool {
	return t.t.Failed()
}

func (t *test) Fatal(args ...any) {
	fmt.Print(t.failedMessage(args...))
	t.exit(1)
}

func (t *test) Helper() {
	t.t.Helper()
}

func (t *test) Setenv(key, value string) {
	t.t.Setenv(key, value)
}

func (t *test) Skip(args ...any) {
	t.t.Skip(args...)
}

func (t *test) Skipped() bool {
	return t.t.Skipped()
}

func (t *test) TempDir() string {
	dir := t.t.TempDir()
	return dir
}

var jsonMarshalIndent = json.MarshalIndent

func (t *test) failedMessage(args ...any) string {
	type value struct {
		Value any    `json:"value"`
		Type  string `json:"type"`
	}
	var message struct {
		Failed   string `json:"failed,omitempty"`
		Result   value  `json:"result,omitempty"`
		Expected value  `json:"expected,omitempty"`
		Message  any    `json:"message,omitempty"`
	}
	message.Failed = t.Caller
	if len(t.result) > 0 {
		message.Result = value{
			Value: t.result,
			Type:  typeOf(t.originalResult),
		}
	}
	if len(t.expected) > 0 {
		message.Expected = value{
			Value: t.expected,
			Type:  typeOf(t.originalExpected),
		}
	}
	if len(args) > 0 {
		message.Message = args[0]
	}

	jsonData, err := jsonMarshalIndent(message, "", "  ")
	if err != nil {
		return ""
	}
	return "\n" + string(jsonData) + "\n\n"
}

func initTest(t T) T {
	ctx := context.Background()
	test := newTest(ctx, t)
	if s, ok := t.(*spy); ok {
		return s
	}
	return test
}

func configureTest(tester T, result, expected any) {
	if t, ok := tester.(*test); ok {
		t.originalResult = result
		t.originalExpected = expected
		t.result = fmt.Sprint(result)
		t.expected = fmt.Sprint(expected)
		t.setCaller(call.Caller())
	}
}
