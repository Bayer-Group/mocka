# mocka <!-- omit in toc -->

The package mocka provides a simple mocking and stubbing library to assist with writing unit tests in go.

All changes will be reflected in the [CHANGELOG](https://github.com/MonsantoCo/mocka/blob/master/CHANGELOG.md)

Currently if you would want to stub a function in go it would be akin to

```go
// alias function for unit testing (in production code)
var jsonMarshal = json.Marshal

...

// create temporary variable to store original function (in unit test)
var jsonMarshalOriginal func(v interface{}) ([]byte, error)

...

jsonMarshalOriginal = jsonMarshal
jsonMarshal = func(v interface{}) ([]byte, error) {
    return []byte("value"), nil
}
defer func() {
    jsonMarshal	= jsonMarshalOriginal
}()
```

This structure increases the length of unit tests; depending on how many functions are needing to be stubbed. Mocka provides a safe way to stub functions while also reducing the amount of code needed.

For example

```go
// alias function for unit testing (in production code)
var jsonMarshal = json.Marshal

...

// create stub to validate against
stub := mocka.Function(&jsonMarshal, []byte("value"), nil)
defer stub.Restore()
```

## Table of Contents <!-- omit in toc -->

- [Public API](#public-api)
- [Stubbing Functions](#stubbing-functions)
	- [Restoring a function's original functionality](#restoring-a-functions-original-functionality)
	- [Changing the return values of a Stub](#changing-the-return-values-of-a-stub)
	- [Changing the return values of a stub based on the call index](#changing-the-return-values-of-a-stub-based-on-the-call-index)
	- [Changing the return values of a Stub based on the arguments](#changing-the-return-values-of-a-stub-based-on-the-arguments)
		- [Changing the return values for a Stub based on the call index of specific arguments](#changing-the-return-values-for-a-stub-based-on-the-call-index-of-specific-arguments)
		- [Changing the return values for a Stub based on argument matchers](#changing-the-return-values-for-a-stub-based-on-argument-matchers)
	- [Retrieving the arguments and return values from a Stub](#retrieving-the-arguments-and-return-values-from-a-stub)
		- [Retrieve how many times the function was called](#retrieve-how-many-times-the-function-was-called)
		- [Retrieve the arguments and return values for all calls against the original function](#retrieve-the-arguments-and-return-values-for-all-calls-against-the-original-function)
		- [Retrieve the arguments and return values for a specific call to the original function](#retrieve-the-arguments-and-return-values-for-a-specific-call-to-the-original-function)
	- [Executing a function when a stub is called](#executing-a-function-when-a-stub-is-called)
- [Creating Sandboxes](#creating-sandboxes)
	- [Stubbing function](#stubbing-function)
	- [Restoring sandboxes](#restoring-sandboxes)

## Public API

The public API for mocka is documented by the following interfaces:

```go
// Call describes the meta data associated to a call of the original function
type Call interface {
    Arguments() []interface{}
    ReturnValues() []interface{}
}

// Returner describes the functionality to update the
// return values for the Stub
type Returner interface {
    Return(...interface{}) error
}

// OnCaller describes the functionality to set custom
// return values based on call index
type OnCaller interface {
    OnCall(int) Returner
    OnFirstCall() Returner
    OnSecondCall() Returner
    OnThirdCall() Returner
}

// GetCaller describes the functionality to get information
// for a specific call to  of the original function
type GetCaller interface {
    GetCalls() []Call
    GetCall(int) Call
    GetFirstCall() Call
    GetSecondCall() Call
    GetThirdCall() Call
    CallCount() int
    CalledOnce() bool
    CalledTwice() bool
    CalledThrice() bool
}

// OnCallReturner describes the functionality to update the
// return values itself of based on the call index of the
// original function
type OnCallReturner interface {
    OnCaller
    Returner
}

// Sandbox describes an isolated environment in which
// functions can be stubbed
type Sandbox interface {
    StubFunction(interface{}, ...interface{}) (Stub, error)
    Restore()
}

// Stub describes the functionality related to the stub
// replacement of a function
type Stub interface {
    Returner
    GetCaller
    OnCaller
    Restore()
    WithArgs(...interface{}) OnCallReturner
    ExecOnCall(func([]interface{}))
}
```

## Stubbing Functions

`mocka.Function` replaces the provided function with a stubbed implementation. The `Stub` has the ability to change the return values of the original function in many different cases. It also provides the ability to get metadata associated to any call against the original function.

`mocka.Function` also returns an error if the replacement of the original function failed.

For example

```go
var fn = func(str string) int {
    return len(str)
}

stub, err := mocka.Function(&fn, 20)
if err != nil {
    log.Fatal(err.Error())
}
defer stub.Restore()

fmt.Println(fn("1"))
// Output: 20
```

### Restoring a function's original functionality

After creating a `Stub` it is recommended to `defer` the `Stub`s recovery. This is to ensure that the `Stub` returns the original functionality back to the function. To restore a `Stub` call the `Restore` function.

For example:

```go
var fn = func(str string) int {
    return len(str)
}

stub, err := mocka.Function(&fn, 20)
if err != nil {
    log.Fatal(err.Error())
}
defer stub.Restore()
```

### Changing the return values of a Stub

mocka allows for the return values of a `Stub` to be changed at any time and in many different cases. When creating `Stub` it is required to specify a default set of return values it will return. If you want to change the default return values after the stub has been created simply call `Return` on the `Stub`.

For example:

```go
var fn = func(str string) int {
    return len(str)
}

stub, err := mocka.Function(&fn, 20)
if err != nil {
    log.Fatal(err.Error())
}
defer stub.Restore()

fmt.Println(fn("123"))

if err = stub.Return(5); err != nil {
    log.Fatal(err.Error())
}

fmt.Println(fn("123"))
// Output: 20
// 5
```

### Changing the return values of a stub based on the call index

mocka allows for return values to be changed based on how many times the original function has been called. To change the return values use the `OnCall` method that can be used by either the `Stub` or a custom set of arguments.

> The _callIndex_ uses zero-based indexing.

For example

```go
var fn = func(str string) int {
    return len(str)
}

stub, err := mocka.Function(&fn, 20)
if err != nil {
    log.Fatal(err.Error())
}
defer stub.Restore()

if err = stub.OnCall(1).Return(5); err != nil {
    log.Fatal(err.Error())
}

if err = stub.OnCall(3).Return(3); err != nil {
    log.Fatal(err.Error())
}

fmt.Println(fn("asdf"))
fmt.Println(fn("234"))
fmt.Println(fn("hello"))
fmt.Println(fn("world"))
// Output: 20
// 5
// 20
// 3
```

mocka provides helper functions for accessing the first three times a function has been called. Instead of using the `OnCall` method the following methods can be used `OnFirstCall`, `OnSecondCall`, or `OnThirdCall`.

### Changing the return values of a Stub based on the arguments

mocka allows for return values to be changed based on the arguments provided to the function. This can be done by using the `WithArgs` method on the `Stub`. `WithArgs` returns the following interface:

```go
type OnCallReturner interface {
	OnCaller
	Returner
}
```

For example

```go
var fn = func(str string) int {
    return len(str)
}

stub, err := mocka.Function(&fn, 20)
if err != nil {
    log.Fatal(err.Error())
}
defer stub.Restore()

if err = stub.WithArgs("123").Return(5); err != nil {
    log.Fatal(err.Error())
}

fmt.Println(fn("123"))
// Output: 5
```

> If `Return` is not called on the `OnCallReturner` interface then it be ignored until `Return` is called.

#### Changing the return values for a Stub based on the call index of specific arguments

Similar to the `Stub` the return values can be changed based on the call index of the original function for a specifc set of arguments. To change the return values for a specific call index use the `OnCall` method.

For example

```go
var fn = func(str string) int {
    return len(str)
}

stub, err := mocka.Function(&fn, 20)
if err != nil {
    log.Fatal(err.Error())
}
defer stub.Restore()

withArgs123 := stub.WithArgs("123")

if err = withArgs123.OnCall(1).Return(5); err != nil {
    log.Fatal(err.Error())
}

if err = withArgs123.OnCall(3).Return(3); err != nil {
    log.Fatal(err.Error())
}

fmt.Println(fn("123"))
fmt.Println(fn("123"))
fmt.Println(fn("123"))
fmt.Println(fn("123"))
// Output: 20
// 5
// 20
// 3
```

mocka provides helper functions for accessing the first three times a function has been called. Instead of using the `OnCall` method the following methods can be used `OnFirstCall`, `OnSecondCall`, or `OnThirdCall`.

#### Changing the return values for a Stub based on argument matchers

mocka provides a powerful `match` package that can be used in conjunction with the `WithArgs` function. Sometimes you might not know the exact value a function is called with. This is a scenario where matchers can help navigate around that problem.

Currently there are over 25 built in matchers you can use. More information can be found at [matcher descriptions](./MATCH.MD).

For example

```go
var fn = func(ctx context.Context, str []string, n int) int {
	// ...
	return 0
}

stub, err := mocka.Function(&fn, 20)
if err != nil {
	log.Fatal(err.Error())
}
defer stub.Restore()

if err = stub.WithArgs(match.Anything(), 2).Return(10); err != nil {
	log.Fatal(err.Error())
}

fmt.Println(fn(context.TODO(), 5))
fmt.Println(fn(context.Background(), 2))
// Output: 20
// 10
```

> The mocka/match package also provides the ability to create your own custom matchers.

### Retrieving the arguments and return values from a Stub

Setting the return values is only half of what mocka can do. Once a `Stub` has been called you can retrieve the arguments and return values the original function was called with.

#### Retrieve how many times the function was called

You can get how many times the original function was called after stubbing the function by using `CallCount`.

For example

```go
var fn = func(str string) int {
    return len(str)
}

stub, err := mocka.Function(&fn, 20)
if err != nil {
    log.Fatal(err.Error())
}
defer stub.Restore()

fn("first call")
fn("second call")
fn("third call")

fmt.Println(stub.CallCount())
// Output: 3
```

mocka provides helper functions for checking if a `Stub` has been called at least the first three times. Instead of using the `CallCount` method the following methods can be used `CalledOnce`, `CalledTwice`, or `CalledThrice`.

#### Retrieve the arguments and return values for all calls against the original function

`GetCalls` returns all calls made to the original function that where captured by the stubbed implementation.

For example

```go
var fn = func(str string) int {
    return len(str)
}

stub, err := mocka.Function(&fn, 20)
if err != nil {
    log.Fatal(err.Error())
}
defer stub.Restore()

fn("first call")
fn("second call")
fn("third call")

calls := stub.GetCalls()
for _, call := range calls {
    fmt.Printf("Arguments: %+v; Return Values: %+v\n", call.Arguments(), call.ReturnValues())
}
// Output: Arguments: [first call]; Return Values: [20]
// Arguments: [second call]; Return Values: [20]
// Arguments: [third call]; Return Values: [20]
```

#### Retrieve the arguments and return values for a specific call to the original function

`GetCall` returns the arguments and return values of the original function that was captured by the stubbed implementation. It will return these values for the specified time the function was called.

`GetCall` will also panic if the call index is lower than zero or greater than the number of times the function was called.

> The call index uses zero-based indexing

For example

```go
var fn = func(str string) int {
    return len(str)
}

stub, err := mocka.Function(&fn, 20)
if err != nil {
    log.Fatal(err.Error())
}
defer stub.Restore()

fn("first call")
fn("second call")
fn("third call")

call := stub.GetCall(2)
fmt.Printf("Arguments: %+v; Return Values: %+v\n", call.Arguments(), call.ReturnValues())
// Output: Arguments: [third call]; Return Values: [20]
```

mocka provides helper functions for retrieving the arguments and return values for the first three calls. Instead of using the `GetCall` method the following methods can be used `GetFirstCall`, `GetSecondCall`, or `GetThirdCall`.

### Executing a function when a stub is called

In some special cases code will need to be run when the original function is called. This code is usually for performing side-effects. Mocka provides the ability to give a `Stub` a function to be called when the original function is called. Call `ExecOnCall` providing a function with the following signature `func(arguments []interface{}) {}` to have it be called when the original function is called. This function will be called with the same arguments the original function is called with.

For example

```go
var fn = func(str string) int {
    return len(str)
}

stub, err := mocka.Function(&fn, 20)
if err != nil {
    log.Fatal(err.Error())
}
defer stub.Restore()

stub.ExecOnCall(func(args []interface{}) {
    fmt.Println(args)
})

fn("123")
// Output: [123]
```

## Creating Sandboxes

`CreateSandbox` returns an isolated sandbox from which functions can be stubbed. The benefit you receive from using a sandbox is the ability to perform one call to `Restore` for a collection of `Stub`s.

### Stubbing function

`StubFunction` replaces the provided function with a stubbed implementation. The stub has the ability to change change the return values of the original function in many different cases. The stub also provides the ability to get metadata associated to any call against the original function.

`StubFunction` also returns an error if the replacement of the original function with the stub failed.

For example

```go
var fn = func(str string) int {
    return len(str)
}

sandbox := mocka.CreateSandbox()
defer sandbox.Restore()

sandbox.StubFunction(&fn, 20)

fmt.Println(fn("1"))
// Output: 20
```

### Restoring sandboxes

After creating a `Sandbox` it is recommended to _defer_ the sandboxes recovery. This is to ensure that the `Stub`s are returned the original functionality. To restore a `Sandbox` call the `Restore` function.

For example:

```go
var fn = func(str string) int {
    return len(str)
}

sandbox := mocka.CreateSandbox()
defer sandbox.Restore()
```