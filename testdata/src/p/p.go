package p

func funcReturnsInt() int {
	return 1
}

func funcReturnsFunc() func() {
	return func() {
	}
}

func funcReturnsErrAndFunc() (error, func()) {
	return nil, func() {
	}
}

func funcDeferReturnGoodValue() {
	defer funcReturnsInt()
}

func funcDeferReturnFunc() {
	defer funcReturnsFunc() // want "deferred call should not return a function"
}

func funcDeferReturnErrAndFunc() {
	defer funcReturnsErrAndFunc() // want "deferred call should not return a function"
}

func funcDeferAnonymousReturnFunc() {
	defer func() func() { // want "deferred call should not return a function"
		return func() {}
	}()
}

func funcDeferAnonymousReturnErrAndFunc() {
	defer func() (error, func()) { // want "deferred call should not return a function"
		return nil, func() {}
	}()
}
