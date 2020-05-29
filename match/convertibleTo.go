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
