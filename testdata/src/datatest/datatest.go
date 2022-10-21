package datatest

type SomeStruct struct {
}

func (s SomeStruct) RetErr() error {
	return nil
}

func (s SomeStruct) RetInAndErr() (int, error) {
	return 100, nil
}

func (s SomeStruct) RetNothing() {
}
