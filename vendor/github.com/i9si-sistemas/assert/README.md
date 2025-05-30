# Assert

The `assert` package provides a set of helper functions for verifying conditions in your tests. It can be used in testing scenarios to assert expected outcomes and handle errors effectively. Each function takes a test context and a message to report any failures.

## Installation

To install the package, run:

```sh
go get github.com/i9si-sistemas/assert
```

## Usage

Import the package in your test files:

```go
import "github.com/i9si-sistemas/assert"
```

### Functions

- **Equal**: Checks if the provided result is equal to the expected value.
  ```go
  assert.Equal(t, result, expected, "Expected result to be equal")
  ```

- **NotEqual**: Checks if the two provided values are not equal.
  ```go
  assert.NotEqual(t, v1, v2, "Values should not be equal")
  ```

- **Nil**: Checks if the provided value is nil.
  ```go
  assert.Nil(t, myVar, "myVar should be nil")
  ```

- **NotNil**: Checks if the provided value is not nil.
  ```go
  assert.NotNil(t, myVar, "myVar should not be nil")
  ```

- **Error**: Checks if the provided error is not nil.
  ```go
  assert.Error(t, err, "Expected an error, but got nil")
  ```

- **NoError**: Checks if the provided error is nil.
  ```go
  assert.NoError(t, err, "Expected no error, but got one")
  ```

- **True**: Checks if the provided boolean value is true.
  ```go
  assert.True(t, result, "Expected result to be true")
  ```

- **False**: Checks if the provided boolean value is false.
  ```go
  assert.False(t, result, "Expected result to be false")
  ```

- **Zero**: Checks if the provided value is equal to the zero value of its type.
  ```go
  assert.Zero(t, value, "Value should be zero")
  ```

### Example

Here is an example of how to use the `assert` package in your tests:

```go
package mypackage

import (
  "testing"
  "github.com/i9si-sistemas/assert"
)

func TestMyFunction(t *testing.T) {
  result := MyFunction()
  expected := 42
  assert.Equal(t, result, expected, "Expected result to be 42")
}
```

## Contributing

We welcome contributions to the `assert` package. If you would like to contribute, please fork the repository and submit a pull request.