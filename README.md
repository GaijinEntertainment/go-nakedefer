<div align="center">

# defer

</div>

---

`defer` is a golang analyzer that finds defer functions which return a function.

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

// valid
func someFuncWithValidDefer1() {
	defer func() {
	}()
}

// valid
func someFuncWithValidDefer2() {
	defer func() error {
		return nil
	}()
}

// invalid, deferred call should not return a function
func someFuncWithInvalidDefer() {
	defer func() func() {
		return func() {

		}
	}()
}
```