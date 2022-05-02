# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [v2.0.1] - 2022-05-02
## Changed
- Updated godoc reference in README.md to point to v2
- change module path to reflect change to Bayer-Group org
- minimum required Go version to 1.16

## [v2.0.0] - 2020-07-27
## Added
- New `SliceOf()` matcher that allows matchers for specific slice elements
- Examples for variadic functions
- Documentation to the README.md for variadic functions
- Migration guide for **v1** to **v2**

## Changed
- `WithArgs()` has been updated to take in variadic arguments the same as how functions take them 
- `argumentValidationError` messages now show variadic arguments as `...type` instead of `[]type`
- `mocka.Function` now takes in a `TestReporter` to fail tests internally inside of mocka
- `mocka.CreateSandbox` now takes in a `TestReporter` to fail tests internally inside of mocka
- `Sandbox.StubFunction` was renamed to `Sandbox.Function` for consistency
- `mocka.Function` and `mocka.CreateSandbox` will call `log.Fatal` if a `TestReporter` is not provided
- Exposed internal types to remove the need to unnecessary interfaces
- `Stub.Return` no longers returns an error, but instead fails the test internally
- `Stub.GetCall` no longer panics, but instead fails the test internally
- Updated documentation to reflect public API changes
- Module is now updated to have v2 suffix
- All internal imports reference mocka/v2

## Removed
- Removed mocka.File()
- `argumentValidationError`
- `outParameterValidationError`
- All public interfaces that were only used for defining the API

## [v1.1.0] - 2020-06-02
## Added
- Revive for linting Go code
- Dlv for debugging Go code
- Additional nil checks to prevent internal panics
- New sub-package `match` containing custom matchers to use with `stub.WithArgs`
- New interface `SupportedKindsMatcher` which custom matchers can be created from
- The following built-in matchers:
  - `Anything()` — matches any value
  - `AnythingButNil()` — matches any value but `nil`
  - `ConvertibleTo(interface{})` — matches if the actual value can be converted to the type of the provided value
  - `ElementsContaining(...interface{})` — matches if the array or slice contains all the provided elements
  - `Empty()` — matches if the string, array, slice, or map have a length of 0
  - `Exactly(interface{})` — matches if the actual value is deep equal to the provided value
  - `KeysContaining(...interface{})` — matches if the actual map contains all the provided keys
  - `LengthOf(int)` — matches if the array, slice, string, or map has the provided length
  - `Nil()` — matches if the value is nil
  - `StringContaining(string)` — matches if the string is containing the provided sub string
  - `StringPrefix(string)` — matches if the string starts with the provided prefix
  - `StringSuffix(string)` — matches if the string ends with the provided suffix
  - `TypeOf(string)` — matches if the actual is of the provided type as a string
  - `ImplementerOf(interface{})` — matches if the actual type implements the interface of the provided value
  - `IntGreaterThan(int64)` — matches if the actual int64 is greater than the provided int64
  - `IntGreaterThanOrEqualTo(int64)` — matches if the actual int64 is greater than or equal to the provided int64
  - `IntLessThan(int64)` — matches if the actual int64 is less than the provided int64
  - `IntLessThanOrEqualTo(int64)` — matches if the actual int64 is less than or equal too the provided int64
  - `UintGreaterThan(uint64)` — matches if the actual uint64 is greater than the provided uint64
  - `UintGreaterThanOrEqualTo(uint64)` — matches if the actual uint64 is greater than or equal to the provided uint64
  - `UintLessThan(uint64)` — matches if the actual uint64 is less than the provided uint64
  - `UintLessThanOrEqualTo(uint64)` — matches if the actual uint64 is less than or equal too the provided uint64
  - `FloatGreaterThan(float32)` — matches if the actual float is greater than the provided float
  - `FloatGreaterThanOrEqualTo(float32)` — matches if the actual float is greater than or equal to the provided float
  - `FloatLessThan(float32)` — matches if the actual float is less than the provided float
  - `FloatLessThanOrEqualTo(float32)` — matches if the actual float is less than or equal too the provided float
- Constructor for the `customArguments` struct
- Testable examples for each mocka matcher
- New function `toFriendlyName` which will return a human readable version of the type

## Changed
- Examples to have more use-cases
- `customArguments.args` -> `customArguments.argMatchers`
- `customArguments` struct to use `[]match.SupportedKindsMatcher` instead of an array of `[]interface{}`
- README to document new and custom matchers

## [v1.0.1] - 2020-01-13
### Added
- Revive for linting Go code
- Dlv for debugging Go code
- Additional nil checks to prevent internal panics
- Mutex to sandbox
- RWMutex to mockFunction

### Changed
- Examples to have more use cases
- Sandbox to use pointer to stubs to not copy mutex value

### Fixed
- Equality checks in WithArgs to use reflect.DeepEqual to compare slices
- Stringified type names to not be empty strings for none primitives

## [v1.0.0] - 2019-07-01
### Added
- Initial public release of this package
