<div align="center">

# defer

</div>

---

`nackedefer` is a golang analyzer that finds defer functions which return anything.

### Installation

```shell
go get -u github.com/GaijinEntertainment/go-nackedefer/cmd/nackedefer
```

### Usage

```
nackedefer [-flag] [package]

Flags:
  -e value
        Regular expression to exclude function names
```


### Example

```go

func funcNotReturnAnyType() {
}

func funcReturnErr() error {
    return errors.New("some error")
}

func funcReturnFuncAndErr() (func(), error) {
    return func() {
    }, nil
}

func ignoreFunc() error {
    return errors.New("some error")
}

func testCaseValid1() {
    defer funcNotReturnAnyType() // valid

    defer func() { // valid
        funcNotReturnAnyType()
    }()

    defer func() { // valid
        _ = funcReturnErr()
    }()
}

func testCaseInvalid1() {
    defer funcReturnErr() // invalid

    defer funcReturnFuncAndErr() // invalid

    defer func() error { // invalid
        return nil
    }()

    defer func() func() { // invalid
        return func() {}
    }()
}

func testCase1() {
    defer fmt.Errorf("some text") // invalid

	r := new(bytes.Buffer)
    defer io.LimitReader(r, 1) // invalid

    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("DONE"))
    }))
    defer srv.Close()                  // invalid
    defer srv.CloseClientConnections() // invalid
    defer srv.Certificate()            // invalid
}

func testCase2() {
    s := datatest.SomeStruct{}
    defer s.RetNothing()  // valid
    defer s.RetErr()      // invalid
    defer s.RetInAndErr() // invalid
}

func testCaseExclude1() {
    // exclude ignoreFunc
    defer ignoreFunc() // valid - excluded 
}

func testCaseExclude2() {
    // exclude os\.(Create|WriteFile|Chmod)
    defer os.Create("file_test1")                                   // valid - excluded 
    defer os.WriteFile("file_test2", []byte("data"), os.ModeAppend) // valid - excluded 
    defer os.Chmod("file_test3", os.ModeAppend)                     // valid - excluded 
    defer os.FindProcess(100500)                                    // invalid
}

func testCaseExclude3() {
    // exclude fmt\.Print.*
    defer fmt.Println("e1")        // valid - excluded 
    defer fmt.Print("e1")          // valid - excluded 
    defer fmt.Printf("e1")         // valid - excluded 
    defer fmt.Sprintf("some text") // invalid
}

func testCaseExclude4() {
    // exclude io\.Close
    rc, _ := zlib.NewReader(bytes.NewReader([]byte("111")))
    defer rc.Close() // valid - excluded 
}
```