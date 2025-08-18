package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	if len(e.errors) == 0 {
		return "0 errors occured"
	}

	var sb strings.Builder
	if len(e.errors) == 1 {
		sb.WriteString("1 error occurred:\n")
	} else {
		fmt.Fprintf(&sb, "%d errors occurred:\n", len(e.errors))
	}

	for _, err := range e.errors {
		fmt.Fprintf(&sb, "\t* %s\n", err.Error())
	}

	return sb.String()
}

func Append(err error, errs ...error) *MultiError {
	var multiError *MultiError

	// find out what is err in first arg
	multiErrorArg, isMultiError := err.(*MultiError)
	if isMultiError {
		multiError = multiErrorArg
	} else if err != nil {
		multiError = &MultiError{errors: []error{err}}
	} else {
		multiError = &MultiError{}
	}

	for _, err := range errs {
		multiError.errors = append(multiError.errors, err)
	}

	return multiError
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occurred:\n\t* error 1\n\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
