<div align="center">

# defer

</div>

---

`defer` is a golang analyzer that finds defer functions which return anything.

### Installation

```shell
go get -u github.com/GaijinEntertainment/go-defer/cmd/defer
```

### Usage

```
defer ./...

```


### Example

```go

func funcNotReturnAnyType() {
}

func funcReturnErr() error {
	return errors.New("some error")
}

// valid
func someFuncWithValidDefer1() {
	defer func() {
	}()
}

// valid
func someFuncWithValidDefer2() {
	defer funcNotReturnAnyType()
}

// invalid, deferred call should not return any type
func someFuncWithInvalidDefer1() {
    defer func() error {
        return nil
    }()
}

// invalid, deferred call should not return any type
func someFuncWithInvalidDefer2() {
    defer funcReturnErr()
}

// invalid, deferred call should not return any type
func someFuncWithInvalidDefer3() {
	defer func() func() {
		return func() {

		}
	}()
}
```