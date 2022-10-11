package p

import (
	"errors"
)

func funcNotReturnAnyType() {
}

func funcReturnErr() error {
	return errors.New("some error")
}

func funcReturnFuncAndErr() (func(), error) {
	return func() {
	}, nil
}

func funcDeferNotReturnAnyType1() {
	defer funcNotReturnAnyType()
}

func funcDeferNotReturnAnyType2() {
	defer func() {
		_ = funcReturnErr()
	}()
}

func funcDeferReturnErr() {
	defer funcReturnErr() // want "deferred call should not return any type"
}

func funcDeferReturnErrAndFunc() {
	defer funcReturnFuncAndErr() // want "deferred call should not return any type"
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
