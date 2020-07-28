# mocka/match

```go
import "github.com/MonsantoCo/mocka/v2/match"
```

The package match provides a powerful suite of matchers that can be used in conjunction with the `WithArgs` function for a mocka Stub.

A matcher is defined by implementing the following interface.

```go
// SupportedKindsMatcher describes the functionality of a custom argument matcher for mocka
type SupportedKindsMatcher interface {
	// SupportedKinds returns the supported kinds for the matcher
	SupportedKinds() map[reflect.Kind]struct{}

	// Match return true is the match was successful; otherwise false
	Match(interface{}) bool
}
```

<details>
<summary>Example</summary>

```go
type anything struct {
}

// SupportedKinds returns all the kinds the anything matcher supports
func (anything) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Bool:          {},
		reflect.Int:           {},
		reflect.Int8:          {},
		reflect.Int16:         {},
		reflect.Int32:         {},
		reflect.Int64:         {},
		reflect.Uint:          {},
		reflect.Uint8:         {},
		reflect.Uint16:        {},
		reflect.Uint32:        {},
		reflect.Uint64:        {},
		reflect.Uintptr:       {},
		reflect.Float32:       {},
		reflect.Float64:       {},
		reflect.Complex64:     {},
		reflect.Complex128:    {},
		reflect.Array:         {},
		reflect.Chan:          {},
		reflect.Func:          {},
		reflect.Interface:     {},
		reflect.Map:           {},
		reflect.Ptr:           {},
		reflect.Slice:         {},
		reflect.String:        {},
		reflect.Struct:        {},
		reflect.UnsafePointer: {},
	}
}

// Match always returns true
func (anything) Match(_ interface{}) bool {
	return true
}
```
</details>

## Built in Matchers

When working with matchers it is possible to have multiple custom arguments match for a set of values. In these scenarios mocka will use the following priority to pick which matcher will be used.

| Matcher                                                           | Priority |
| ----------------------------------------------------------------- | -------- |
| [Exactly](#exactly)                                               | 25       |
| [Nil](#nil)                                                       | 24       |
| [Float Greater Than](#float-greater-than)                         | 23       |
| [Float Less Than](#float-less-than)                               | 22       |
| [Float Greater Than Or Equal To](#float-greater-than-or-equal-to) | 21       |
| [Float Less Than Or Equal To](#float-less-than-or-equal-to)       | 20       |
| [IntGreaterThan](#int-greater-than)                               | 19       |
| [Int LessThan](#int-less-than)                                    | 18       |
| [Int GreaterThanOrEqualTo](#int-greater-than-or-equal-to)         | 17       |
| [Int LessThanOrEqualTo](#int-less-than-or-equal-to)               | 16       |
| [Uint Greater Than](#uint-greater-than)                           | 15       |
| [Uint Less Than](#uint-less-than)                                 | 14       |
| [Uint Greater Than Or Equal To](#uint-greater-than-or-equal-to)   | 13       |
| [Uint Less Than Or Equal To](#uint-less-than-or-equal-to)         | 12       |
| [String Prefix](#string-prefix)                                   | 11       |
| [String Suffix](#string-suffix)                                   | 10       |
| [String Containing](#string-containing)                           | 9        |
| [Length Of](#length-of)                                           | 8        |
| [Empty](#empty)                                                   | 7        |
| [Keys Containing](#keys-containing)                               | 6        |
| [Elements Containing](#elements-containing)                       | 5        |
| [Implementer Of](#implementer-of)                                 | 4        |
| [Convertible To](#convertible-to)                                 | 3        |
| [Type Of](#type-of)                                               | 2        |
| [Anything But Nil](#anything-but-nil)                             | 1        |
| [Anything](#anything)                                             | 0        |


> If you are using a custom matcher (non built in matcher) it's priority will be the highest priority.


## Exact Value Matchers

### Exactly
---

The `Exactly(interface{})` matcher will match only if the value is deep equal to the provided value.

> `Exactly` is the default matcher for all values provided to `WithArgs` if a matcher is not supplied, unless the value is `nil`.

<details>
<summary>Example</summary>

```go
match.Exactly([]int{1,2,3})
```

</details>

#### Supported Kinds

Bool, Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Uintptr, Float32, Float64, Complex64, Complex128, Array, Chan, Func, Interface, Map, Ptr, Slice, String, Struct, UnsafePointer

### Nil
---

The `Nil()` matcher will match any nil value.

> `Nil` is the default matcher for `nil` values provided to `WithArgs`.

#### Supported Kinds

Chan, Func, Interface, Map, Ptr, Slice

## Numeric Matchers

### Float Greater Than
---

The `FloatGreaterThan(float64)` matcher will match if the numeric value is greater than the provided value.

<details>
<summary>Example</summary>

```go
match.FloatGreaterThan(2)
```

</details>

#### Supported Kinds

Float32, Float64

### Float Less Than
---

The `FloatLessThan(float64)` matcher will match if the numeric value is less than the provided value.

<details>
<summary>Example</summary>

```go
match.FloatLessThan(2)
```

</details>

#### Supported Kinds

Float32, Float64

### Float Greater Than Or Equal To
---

The `FloatGreaterThanOrEqualTo(float64)` matcher will match if the numeric value is greater than or equal to the provided value.

<details>
<summary>Example</summary>

```go
match.FloatGreaterThanOrEqualTo(2)
```

</details>

#### Supported Kinds

Float32, Float64

### Float Less Than Or Equal To
---

The `FloatLessThanOrEqualTo(float64)` matcher will match if the numeric value is less than or equal to the provided value.

<details>
<summary>Example</summary>

```go
match.FloatLessThanOrEqualTo(2)
```

</details>

#### Supported Kinds

Float32, Float64

### Int Greater Than
---

The `IntGreaterThan(int64)` matcher will match if the numeric value is greater than the provided value.

<details>
<summary>Example</summary>

```go
match.IntGreaterThan(2)
```

</details>

#### Supported Kinds

Int, Int8, Int16, Int32, Int64

### Int Less Than
---

The `IntLessThan(int64)` matcher will match if the numeric value is less than the provided value.

<details>
<summary>Example</summary>

```go
match.IntLessThan(2)
```

</details>

#### Supported Kinds

Int, Int8, Int16, Int32, Int64

### Int Greater Than Or Equal To
---

The `IntGreaterThanOrEqualTo(int64)` matcher will match if the numeric value is greater than or equal to the provided value.

<details>
<summary>Example</summary>

```go
match.IntGreaterThanOrEqualTo(2)
```

</details>

#### Supported Kinds

Int, Int8, Int16, Int32, Int64

### Int Less Than Or Equal To
---

The `IntLessThanOrEqualTo(int64)` matcher will match if the numeric value is less than or equal to the provided value.

<details>
<summary>Example</summary>

```go
match.IntLessThanOrEqualTo(2)
```

</details>

#### Supported Kinds

Int, Int8, Int16, Int32, Int64

### Uint Greater Than
---

The `UintGreaterThan(uint64)` matcher will match if the numeric value is greater than the provided value.

<details>
<summary>Example</summary>

```go
match.UintGreaterThan(2)
```

</details>

#### Supported Kinds

Uint, Uint8, Uint16, Uint32, Uint64

### Uint Less Than
---

The `UintLessThan(uint64)` matcher will match if the numeric value is less than the provided value.

<details>
<summary>Example</summary>

```go
match.UintLessThan(2)
```

</details>

#### Supported Kinds

Uint, Uint8, Uint16, Uint32, Uint64

### Uint Greater Than Or Equal To
---

The `UintGreaterThanOrEqualTo(uint64)` matcher will match if the numeric value is greater than or equal to the provided value.

<details>
<summary>Example</summary>

```go
match.UintGreaterThanOrEqualTo(2)
```

</details>

#### Supported Kinds

Uint, Uint8, Uint16, Uint32, Uint64

### Uint Less Than Or Equal To
---

The `UintLessThanOrEqualTo(uint64)` matcher will match if the numeric value is less than or equal to the provided value.

<details>
<summary>Example</summary>

```go
match.UintLessThanOrEqualTo(2)
```

</details>

#### Supported Kinds

Uint, Uint8, Uint16, Uint32, Uint64

## String Matchers

### String Prefix
---

The `StringPrefix(string)` matcher will match a value if the string starts with the provided string.

<details>
<summary>Example</summary>

```go
match.StringPrefix("he")
```

</details>

#### Supported Kinds

String

### String Suffix
---

The `StringSuffix(string)` matcher will match a value if the string ends with the provided string.

<details>
<summary>Example</summary>

```go
match.StringSuffix("ello")
```

</details>

#### Supported Kinds

String

### String Containing
---

The `StringContaining(string)` matcher will match a value if the string contains the provided string.

<details>
<summary>Example</summary>

```go
match.StringContaining("ello")
```

</details>

#### Supported Kinds

String

## Multiple Purpose Matchers

### Length Of
---

The `LengthOf(int)` matcher will match a value if it's length matches the provided length.

<details>
<summary>Example</summary>

```go
match.LengthOf(12)
```

</details>

#### Supported Kinds

Array, Map, Slice, String

### Empty
---

The `Empty()` matcher will match a value if it has a length of 0.

#### Supported Kinds

Array, Map, Slice, String

### Keys Containing
---

The `KeysContaining(...interface{})` matcher will match a value if all elements exists are keys in the provided map.
  
<details>
<summary>Example</summary>

```go
match.KeysContaining("A", "B")
```

</details>

#### Supported Kinds

Map

### Elements Containing
---

The `ElementsContaining(...interface{})` matcher will match a value if all elements are contained within the provided value.

<details>
<summary>Example</summary>

```go
match.ElementsContaining("A", "B")
```

</details>

#### Supported Kinds

Array, Slice

## Type Matchers

### Implementer Of
---

The `ImplementerOf(interface{})` matcher will match a value if it is a pointer that implements the provided interface type.

<details>
<summary>Example</summary>

```go
match.ImplementerOf(new(error))
```

</details>

#### Supported Kinds

Ptr

### Convertible To
---

The `ConvertibleTo(interface{})` matcher will match a value if it's type can be converted to the provided type.

<details>
<summary>Example</summary>

```go
match.ConvertibleTo(new(int64))
```

</details>

#### Supported Kinds

Bool, Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Uintptr, Float32, Float64, Complex64, Complex128, Array, Chan, Func, Interface, Map, Ptr, Slice, String, Struct, UnsafePointer

### Type Of
---

The `TypeOf(string)` matcher will match a value if it's type is the same of the provided string.

<details>
<summary>Example</summary>

```go
match.TypeOf("int")
```

</details>

#### Supported Kinds

Bool, Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Uintptr, Float32, Float64, Complex64, Complex128, Array, Chan, Func, Interface, Map, Ptr, Slice, String, Struct, UnsafePointer

### Anything But Nil
---

The `AnythingButNil()` matcher will match any value except `nil`.

#### Supported Kinds

Chan, Func, Interface, Map, Ptr, Slice

### Anything
---

The `Anything()` matcher will match any value regardless of the value provided

#### Supported Kinds

Bool, Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Uintptr, Float32, Float64, Complex64, Complex128, Array, Chan, Func, Interface, Map, Ptr, Slice, String, Struct, UnsafePointer
