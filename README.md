# mocka [![Build Status][build-badge]][build-ci] [![gopherbadger-tag-do-not-edit][coverage-badge]][coverage] [![GoDoc][godoc-badge]][godoc]

```go
import "github.com/Bayer-Group/mocka/v2"
```

Mocka is a simple mocking and stubbing library for the [Go programming language][golang]. It is used to assist with writing unit tests around third-party functions.

All changes will be reflected in the [CHANGELOG][changelog].

> If you are looking to migrate from **v1** to **v2** check out the [migration guide][migrationGuide].

## Why Mocka?

There are times when you would want to control the output of a third-party function in testing. Sometimes making a wrapper around that package/function is more effort than it is worth. Mocka is here to solve that problem. It allows you to control the output of functions without needing to write any additional code. 

Currently if you would want to control the output of a function in go it would be akin to

```go
// --- main.go ---

// alias function for unit testing
var jsonMarshal = json.Marshal

...

// --- main_test.go ---

// create temporary variable to store original function
var jsonMarshalOriginal func(v interface{}) ([]byte, error)

func TestMarshal(t *testing.T) {
    jsonMarshalOriginal = jsonMarshal
    jsonMarshal = func(v interface{}) ([]byte, error) {
        return []byte("value"), nil
    }
    defer func() {
        jsonMarshal = jsonMarshalOriginal
    }()
    
    // Your test code
}
```

This structure increases the length of unit tests; depending on how many functions are needing to control. Mocka provides a safe way to stub functions while also reducing the amount of code required.

> Mocka does this safely using reflection, no calls to the `unsafe` package are made.

The mocka way would be

```go
// --- main.go ---

// alias function for unit testing (in production code)
var jsonMarshal = json.Marshal

...

// --- main_test.go ---

func TestMarshal(t *testing.T) {
    stub := mocka.Function(t, &jsonMarshal, []byte("value"), nil)
    defer stub.Restore()
    
    // Your test code
}
```

> The `encoding/json` package was used in examples for simplicity and not for the need to control it's output. 

## Test Reporter

There are some cases when interacting with a stub where errors can occur. Mocka uses a custom interface called `TestReporter`, which is defined below, to fail tests for you.

```go
type TestReporter interface {
    Errorf(string, ...interface{})
}
```

`TestReporter` is satisfied by the built-in `testing.T` and other testing frameworks like [Ginkgo][ginkgo] by using `GinkgoT()`. 

## Stubs

### Creating a Stub

```go
func Function(
    testReporter TestReporter,
    functionPointer interface{},
    returnValues ...interface{}) *Stub {
        
    }
```

`mocka.Function` replaces the provided function with a stubbed implementation. The `Stub` has the ability to change the return values of the original function in many different cases. It also provides the ability to get metadata associated to any call against the original function.

### Restoring a function's original functionality

After creating a `Stub` it is recommended to `defer` it's restoration. This is to ensure that the `Stub` returns the original functionality back to the function. To restore a `Stub` call the `Restore` function.

<details>
<summary>Example</summary>

```go
package main

import (
    "testing"

    "github.com/Bayer-Group/mocka/v2"
)

func TestMocka(t *testing.T) {
    fn := func(str string) int {
        return len(str)
    }

    stub := mocka.Function(t, &fn, 20)
    defer stub.Restore()

    actual := fn("1")
    if actual != 20 {
        t.Errorf("expected 20 but got %v", actual)
    }
}
```

</details>

### Changing the return values of a Stub

Mocka allows for the return values of a `Stub` to be changed at any time and in many different cases. When creating `Stub` it is required to specify a default set of return values it will return. If you want to change the default return values after the stub has been created simply call `Return` on the `Stub`.

<details>
<summary>Example</summary>

```go
package main

import (
    "testing"

    "github.com/Bayer-Group/mocka/v2"
)

func TestMocka(t *testing.T) {
    fn := func(str string) int {
        return len(str)
    }

    stub := mocka.Function(t, &fn, 20)
    defer stub.Restore()

    if actual := fn("123"); actual != 20 {
        t.Errorf("expected 20 but got %v", actual)
    }

    stub.Return(5)

    if actual := fn("123"); actual != 5 {
        t.Errorf("expected 5 but got %v", actual)
    }
}
```

</details>

### Changing the return values of a stub based on the call index

Mocka allows for return values to be changed based on how many times the original function has been called. To change the return values use the `OnCall` method that can be used by either the `Stub` or a custom set of arguments.

> The _callIndex_ uses zero-based indexing.

Mocka provides helper functions for accessing the first three times a function has been called. Instead of using the `OnCall` method the following methods can be used `OnFirstCall`, `OnSecondCall`, or `OnThirdCall`.

<details>
<summary>Example</summary>

```go
package main

import (
    "testing"

    "github.com/Bayer-Group/mocka/v2"
)

func TestMocka(t *testing.T) {
    fn := func(str string) int {
        return len(str)
    }

    stub := mocka.Function(t, &fn, 20)
    defer stub.Restore()

    withArgs123 := stub.WithArgs("123")

    withArgs123.OnCall(1).Return(5)
    withArgs123.OnCall(3).Return(3)
    
    if actual := fn("123"); actual != 20 {
        t.Errorf("expected 20 but got %v", actual)
    }
    
    if actual := fn("123"); actual != 5 {
        t.Errorf("expected 5 but got %v", actual)
    }
    
    if actual := fn("123"); actual != 20 {
        t.Errorf("expected 20 but got %v", actual)
    }
    
    if actual := fn("123"); actual != 3 {
        t.Errorf("expected 3 but got %v", actual)
    }
}
```

</details>

### Changing the return values of a Stub based on the arguments

Mocka allows for return values to be changed based on the arguments provided to the function. This can be done by using the `WithArgs` method on the `Stub`.

> If `Return` is not called on the `OnCallReturner` interface then it be ignored until `Return` is called.

<details>
<summary>Example</summary>

```go
package main

import (
    "testing"

    "github.com/Bayer-Group/mocka/v2"
)

func TestMocka(t *testing.T) {
    fn := func(str []string, n int) int {
        return len(str) + n
    }

    stub := mocka.Function(t, &fn, 20)
    defer stub.Restore()

    stub.WithArgs([]string{"123", "456"}, 2).Return(5)

    fmt.Println(fn([]string{"123", "456"}, 2))

    if actual := fn([]string{"123", "456"}, 2); actual != 5 {
        t.Errorf("expected 5 but got %v", actual)
    }
}
```

</details>

#### Changing the return values for a Stub based on variadic arguments

mocka accepts variadic arguments for `WithArgs` the same as if you were calling the function itself.

> You can still pass in custom matchers from the `match` package for each element in the variadic list.

<details>
<summary>Example</summary>

```go
package main

import (
    "testing"

    "github.com/Bayer-Group/mocka/v2"
)

func TestMocka(t *testing.T) {
    fn := func(str string, opts ...string) int {
        return len(str) + len(opts)
    }

    stub := mocka.Function(t, &fn, 20)
    stub.WithArgs("A", "B", "C").Return(5)

    if actual := fn("A", "B", "C"); actual != 5 {
        t.Errorf("expected 5 but got %v", actual)
    }
    
    if actual := fn("A"); actual != 20 {
        t.Errorf("expected 20 but got %v", actual)
    }
}
```

</details>

#### Changing the return values for a Stub based on the call index of specific arguments

Similar to the `Stub` the return values can be changed based on the call index of the original function for a specifc set of arguments. To change the return values for a specific call index use the `OnCall` method.

mocka provides helper functions for accessing the first three times a function has been called. Instead of using the `OnCall` method the following methods can be used `OnFirstCall`, `OnSecondCall`, or `OnThirdCall`.

<details>
<summary>Example</summary>

```go
package main

import (
    "testing"

    "github.com/Bayer-Group/mocka/v2"
)

func TestMocka(t *testing.T) {
    fn := func(str string) int {
        return len(str)
    }

    stub := mocka.Function(t, &fn, 20)
    defer stub.Restore()

    withArgs123 := stub.WithArgs("123")

    withArgs123.OnCall(1).Return(5)
    withArgs123.OnCall(3).Return(3)
    
    if actual := fn("123"); actual != 20 {
        t.Errorf("expected 20 but got %v", actual)
    }
    
    if actual := fn("123"); actual != 5 {
        t.Errorf("expected 5 but got %v", actual)
    }
    
    if actual := fn("123"); actual != 20 {
        t.Errorf("expected 20 but got %v", actual)
    }
    
    if actual := fn("123"); actual != 3 {
        t.Errorf("expected 3 but got %v", actual)
    }
}
```

</details>

#### Changing the return values for a Stub based on argument matchers

mocka provides a powerful `match` package that can be used in conjunction with the `WithArgs` function. Sometimes you might not know the exact value a function is called with. This is a scenario where matchers can help navigate around that problem.

Currently there are over 25 built in matchers you can use. More information can be found at [matcher descriptions](./MATCH.md).

> The `match` package also provides the ability to create your own custom matchers.

<details>
<summary>Example</summary>

```go
package main

import (
    "testing"

    "github.com/Bayer-Group/mocka/v2"
)

func TestMocka(t *testing.T) {
    fn := func(str []string, n int) int {
        return len(str) + n
    }

    stub := mocka.Function(t, &fn, 20)
    defer stub.Restore()

    stub.WithArgs(match.Anything(), 2).Return(10)
    stub.WithArgs([]string{"123", "456"}, 2).Return(5)

    if actual := fn([]string{"hello"}, 5); actual != 20 {
        t.Errorf("expected 20 but got %v", actual)
    }
    
    if actual := fn([]string{"mocka"}, 2); actual != 10 {
        t.Errorf("expected 10 but got %v", actual)
    }
    
    if actual := fn([]string{"123", "456"}, 2); actual != 5 {
        t.Errorf("expected 5 but got %v", actual)
    }
}
```

</details>

### Retrieving the arguments and return values from a Stub

Setting the return values is only half of what mocka can do. Once a `Stub` has been called you can retrieve the arguments and return values the original function was called with.

#### Retrieve how many times the function was called

You can get how many times the original function was called after stubbing the function by using `CallCount`.

mocka provides helper functions for checking if a `Stub` has been called at least the first three times. Instead of using the `CallCount` method the following methods can be used `CalledOnce`, `CalledTwice`, or `CalledThrice`.

<details>
<summary>Example</summary>

```go
package main

import (
    "testing"

    "github.com/Bayer-Group/mocka/v2"
)

func TestMocka(t *testing.T) {
    fn := func(str string) int {
        return len(str)
    }

    stub := mocka.Function(t, &fn, 20)
    defer stub.Restore()

    fn("first call")
    fn("second call")
    fn("third call")

    if actual := stub.CallCount(); actual != 3 {
        t.Errorf("expected 3 but got %v", actual)
    }
}
```

</details>

#### Retrieve the arguments and return values for all calls against the original function

`GetCalls` returns all calls made to the original function that where captured by the stubbed implementation.

<details>
<summary>Example</summary>

```go
package main

import (
    "testing"

    "github.com/Bayer-Group/mocka/v2"
)

type test struct {
    arguments []interface{}
    returnValues []interface{}
}

func TestMocka(t *testing.T) {
    fn := func(str string) int {
        return len(str)
    }

    stub := mocka.Function(t, &fn, 20)
    defer stub.Restore()

    fn("first call")
    fn("second call")
    fn("third call")

    calls := stub.GetCalls()
    if len(calls) != 3 {
        t.Fatalf("expected 3 but got %v", actual)
    }
    
    tests := []test{
        {arguments: []interface{}{"first call"}, returnValues: []interface{}{20}},
        {arguments: []interface{}{"second call"}, returnValues: []interface{}{20}},
        {arguments: []interface{}{"third call"}, returnValues: []interface{}{20}},
    }
    
    for i, tc := range tests {
        call := calls[i]
        if !reflect.DeepEqual(tc.arguments, call.Arguments()) {
            t.Fatalf("expected arguments: %v, got: %v", tc.arguments, call.Arguments())
        }
        
        if !reflect.DeepEqual(tc.returnValues, call.ReturnValues()) {
            t.Fatalf("expected return values: %v, got: %v", tc.returnValues, call.ReturnValues())
        }
    }
}
```

</details>

#### Retrieve the arguments and return values for a specific call to the original function

`GetCall` returns the arguments and return values of the original function that was captured by the stubbed implementation. It will return these values for the specified time the function was called.

`GetCall` will also panic if the call index is lower than zero or greater than the number of times the function was called.

> The call index uses zero-based indexing

mocka provides helper functions for retrieving the arguments and return values for the first three calls. Instead of using the `GetCall` method the following methods can be used `GetFirstCall`, `GetSecondCall`, or `GetThirdCall`.

<details>
<summary>Example</summary>

```go
package main

import (
    "testing"

    "github.com/Bayer-Group/mocka/v2"
)

func TestMocka(t *testing.T) {
    fn := func(str string) int {
        return len(str)
    }

    stub := mocka.Function(t, &fn, 20)
    defer stub.Restore()

    fn("first call")
    fn("second call")
    fn("third call")

    call := stub.GetCall(2)
    if !reflect.DeepEqual([]interface{}{"third call"}, call.Arguments()) {
        t.Fatalf("expected arguments: %v, got: %v", []interface{}{"third call"}, call.Arguments())
    }
    
    if !reflect.DeepEqual([]interface{}{20}, call.ReturnValues()) {
        t.Fatalf("expected return values: %v, got: %v", []interface{}{20}, call.ReturnValues())
    }
}
```

</details>


### Executing a function when a stub is called

In some special cases code will need to be run when the original function is called. This code is usually for performing side-effects. Mocka provides the ability to give a `Stub` a function to be called when the original function is called. Call `ExecOnCall` providing a function with the following signature `func(arguments []interface{}) {}` to have it be called when the original function is called. This function will be called with the same arguments the original function is called with.

<details>
<summary>Example</summary>

```go
package main

import (
    "testing"

    "github.com/Bayer-Group/mocka/v2"
)

func TestMocka(t *testing.T) {
    fn := func(in <-chan int) <-chan int {
        out := make(chan int, 1)
        go func() {
            out <- <-in
        }()
        return out
    }

    out := make(chan int, 1)
    stub := mocka.Function(t, &fn, out)
    defer stub.Restore()

    stub.ExecOnCall(func(args []interface{}) {
        c := args[0].(<-chan int)
        out <- <-c
    })

    in := make(chan int, 1)
    in <- 10
    if actual := <-fn(in); actual != 10 {
        t.Fatalf("expected: 10 got: %v", actual)
    }
}
```

</details>

## Sandboxes

In many cases you might need to stub out many functions in a single test file. A `Sandbox` allows you to simplify the restoration of many stubbed functions. You can create one `Sandbox` where you can only call `.Restore()` once for all stubbed functions. 

To create a `Sandbox` call `mocka.CreateSandbox` passing in a [test reporter](#test-reporter). The test reporter will be used to fail any tests where a stubbing error occurs. All stubs created from the `Sandbox` will use the same test reporter.

### API

### Creating a `Sandbox`

```go
func Function(functionPointer interface{}, returnValues ...interface{}) {}
```

`Sandbox.Function` behaves the same as `mocka.Function`. It replaces the provided function with a stubbed implementation. The stub has the ability to change change the return values of the original function in many different cases. The stub also provides the ability to get metadata associated to any call against the original function.

### Restoring a `Sandbox`

```go
func Sandbox.Restore() {}
```

`Sandbox.Restore` will call `.Restore()` on all stubs that have been created from the sandbox. Once the stubs have been restored they are removed from the sandbox. To ensure no other tests are effected by the stubs created from a `Sandbox`, restore it after every test. 

It is recommended to call `Sandbox.Restore` in a _defer_ directly after the sandboxes creation. If you are using a different testing package like [Ginkgo][ginkgo] then placing the restoration call in the `AfterEach(func())` will work as well.


<details>
<summary>Example</summary>

```go
package main

import (
    "testing"

    "github.com/Bayer-Group/mocka/v2"
)

func TestSandbox(t *testing.T) {
    fn := func(str string) int {
        return len(str)
    }

    sandbox := mocka.CreateSandbox(t)
    defer sandbox.Restore()

    sandbox.Function(&fn, 20)

    actual := fn("1")
    if actual != 20 {
        t.Errorf("expected 20 but got %v", actual)
    }
}
```
</details>

[changelog]: https://github.com/Bayer-Group/mocka//blob/master/CHANGELOG.md
[coverage]: https://github.com/jpoles1/gopherbadger
[coverage-badge]: https://img.shields.io/badge/Go%20Coverage-100%25-brightgreen.svg?longCache=true&style=flat
[golang]:          http://golang.org/
[golang-install]:  http://golang.org/doc/install.html#releases
[build-badge]: https://github.com/Bayer-Group/mocka//workflows/build/badge.svg
[build-ci]:       https://github.com/Bayer-Group/mocka//actions?query=workflow%3A%22build%22
[godoc-badge]:     https://godoc.org/github.com/Bayer-Group/mocka/?status.svg
[godoc]:           https://pkg.go.dev/github.com/Bayer-Group/mocka/v2?tab=doc
[ginkgo]: https://github.com/onsi/ginkgo
[migrationGuide]: https://github.com/Bayer-Group/mocka//blob/master/MIGRATE_TO_V2.md

