package p

import (
	"errors"
)

func funcNotReturnsAnyType() {
}

func funcReturnsErr() error {
	return errors.New("some error")
}

func funcReturnsFuncAndErr() (func(), error) {
	return func() {
	}, nil
}

func funcDeferNotReturnAnyType1() {
	defer funcNotReturnsAnyType()
}

func funcDeferNotReturnAnyType2() {
	defer func() {
		_ = funcReturnsErr()
	}()
}

func funcDeferReturnErr() {
	defer funcReturnsErr() // want "deferred call should not return any type"
}

func funcDeferReturnErrAndFunc() {
	defer funcReturnsFuncAndErr() // want "deferred call should not return any type"
}

func funcDeferAnonymousReturnFunc() {
	defer func() func() { // want "deferred call should not return any type"
		return func() {}
	}()
}

func funcDeferAnonymousReturnIntAndErr() {
	defer func() (int, error) { // want "deferred call should not return any type"
		return 1, nil
	}()
}
