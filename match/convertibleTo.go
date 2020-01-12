package match

import (
	"reflect"
)

// ConvertibleTo returns a new matcher that will match a value if it can be converted to the provided type
func ConvertibleTo(value interface{}) SupportedKindsMatcher {
	return &convertibleTo{value}
}

type convertibleTo struct {
	value interface{}
}

// SupportedKinds returns all the kinds the convertible to matcher supports
func (m *convertibleTo) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Bool:          struct{}{},
		reflect.Int:           struct{}{},
		reflect.Int8:          struct{}{},
		reflect.Int16:         struct{}{},
		reflect.Int32:         struct{}{},
		reflect.Int64:         struct{}{},
		reflect.Uint:          struct{}{},
		reflect.Uint8:         struct{}{},
		reflect.Uint16:        struct{}{},
		reflect.Uint32:        struct{}{},
		reflect.Uint64:        struct{}{},
		reflect.Uintptr:       struct{}{},
		reflect.Float32:       struct{}{},
		reflect.Float64:       struct{}{},
		reflect.Complex64:     struct{}{},
		reflect.Complex128:    struct{}{},
		reflect.Array:         struct{}{},
		reflect.Chan:          struct{}{},
		reflect.Func:          struct{}{},
		reflect.Interface:     struct{}{},
		reflect.Map:           struct{}{},
		reflect.Ptr:           struct{}{},
		reflect.Slice:         struct{}{},
		reflect.String:        struct{}{},
		reflect.Struct:        struct{}{},
		reflect.UnsafePointer: struct{}{},
	}
}

// Match returns true the value can be converted to the provided interface; otherwise false.
// If provided value is nil the matcher will always return false
func (m *convertibleTo) Match(value interface{}) bool {
	if m.value == nil || value == nil {
		return false
	}

	expectedType := reflect.TypeOf(m.value)
	actualType := reflect.TypeOf(value)

	if expectedType.Kind() != reflect.Ptr {
		return false
	}

	return actualType.ConvertibleTo(expectedType.Elem())
}
