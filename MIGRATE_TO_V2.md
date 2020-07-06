# Migrating to V2.x.x <!-- omit in toc -->

The release of **v2** comes with many changes to the public API for mocka. This guide is to help walk you through those changes in your project. 

## Table of Contents <!-- omit in toc -->
- [Pull in **v2**](#pull-in-v2)
- [The Way to Create Stubs](#the-way-to-create-stubs)
  - [Stubs](#stubs)
  - [Sandboxes](#sandboxes)
- [Remove Error Handling](#remove-error-handling)
- [New Support for Variadic Functions](#new-support-for-variadic-functions)
- [Goodbye Public API Interfaces](#goodbye-public-api-interfaces)
- [Goodbye Panics](#goodbye-panics)
- [Goodbye File Mocks](#goodbye-file-mocks)

## Pull in **v2**
To start pulling in **v2** update your _go.mod_ to point to the new version. 

Change
```go
// x.x is used to denote any minor or patch version for v1
github.com/MonsantoCo/mocka v1.x.x
```

To
```
github.com/MonsantoCo/mocka v2.0.0
```

## The Way to Create Stubs
Mocka now will report errors and fail tests for you. This is an added benefit which takes the burden of handling mocka errors off the caller. In order for mocka to do this however it needs a way to report those failures. `mocka.Function` and `mocka.CreateSandbox` now accept a `mocka.TestReporter` as the first parameter. The built-in `*testing.T` or `GinkgoT()` from [Ginkgo][ginkgo] satisfy this interface.

The following function will be used in migration examples.
```go
fn := func(_ string) int {
    return 0
}
```

### Stubs
For standard stubs remove the error check and add either `*testing.T` or `GinkgoT()` as the first parameters to `mocka.Function`.

Example **v1**
```go
stub, err := mocka.Function(&_fn, 1)
if err != nil {
    t.Errorf("expected nil but got %v", err.Error())
}
defer stub.Restore()
```

Example **v2**
```go
stub := mocka.Function(t, &_fn, 1)
defer stub.Restore()
```

### Sandboxes
Sandboxes are similar to stubs where `mocka.CreateSandbox` now takes in a new parameters for the `mocka.TestReporter`. In addition to this change the `sandbox.StubFunction` method was renamed to `sandbox.Function` to more closely match the stub API.

Example **v1**
```go
sandbox := mocka.CreateSandbox()
defer sandbox.Restore()

stub, err := sandbox.StubFunction(&_fn, 1)
if err != nil {
    t.Errorf("expected nil but got %v", err.Error())
}
```

Example **v2**
```go
sandbox := mocka.CreateSandbox(t)
defer sandbox.Restore()

stub := sandbox.Function(&_fn, 1)
```

> Note: `sandbox.Function` does not take in a `mocka.TestReporter`. Every stub created from the sandbox will use the same reporter that the sandbox was created with.

## Remove Error Handling

Since mocka now handles errors internally you will need to remove them in your project. In **v1** anytime the `stub.Return` method was called it returned an error. 

The following function will be used in migration examples.
```go
fn := func(_ string) int {
    return 0
}
```

Example **v1**
```go
stub, err := mocka.Function(&_fn, 1)
if err != nil {
    t.Errorf("expected nil but got %v", err.Error())
}
defer stub.Restore()

err = stub.Return(2)
if err != nil {
    t.Errorf("expected nil but got %v", err.Error())
}

err = stub.WithArgs("hello").Return(3)
if err != nil {
    t.Errorf("expected nil but got %v", err.Error())
}
```

Example **v2**
```go
stub := mocka.Function(t, &_fn, 1)
defer stub.Restore()

stub.Return(2)
stub.WithArgs("hello").Return(3)
```

## New Support for Variadic Functions

In **v1** mocka treated variadic function arguments as slices. This matched closely to how the arguments are handled in reflection. However, this did not make sense for the public API. Now in **v2** variadic arguments are treated the same in mocka as they are in your production project code.

The following function will be used in migration examples.
```go
fn := func(_ string, _ ...string) int {
    return 0
}
```

Example **v1**
```go
stub, err := mocka.Function(&_fn, 1)
if err != nil {
    t.Errorf("expected nil but got %v", err.Error())
}
defer stub.Restore()

err = stub.WithArgs("hello", []string{"w", "o", "r", "l", "d"}).Return(3)
if err != nil {
    t.Errorf("expected nil but got %v", err.Error())
}
```

Example **v2** _(variadic argument supplied)_
```go
stub := mocka.Function(t, &_fn, 1)
defer stub.Restore()

stub.WithArgs("hello", "w", "o", "r", "l", "d").Return(3)
```

Example **v2** _(variadic argument ignored)_
```go
stub := mocka.Function(t, &_fn, 1)
defer stub.Restore()

stub.WithArgs("hello").Return(3)
```

## Goodbye Public API Interfaces

Any public interfaces that were _only_ needed to define the API have been removed. This should not effect the majority of users. In that case that your project was using these interfaces it is safe to swap to the new publicly exposed types instead.  

## Goodbye Panics

In **v1** `stub.GetCall` would panic if the call index was out of bounds. This has been changed to fail a test using the supplied reporter. You can remove any safe gaurds in your project the caught panics from mocka now in **v2**.

## Goodbye File Mocks

`mocka.File` was deprecated in **v1** and is now removed in **v2**. File mocks have not fully worked since the release of go 1.13.x without extra work from the caller. It is recommended to refactor to use oi interfaces and mocks moving forword.

[ginkgo]: https://github.com/onsi/ginkgo