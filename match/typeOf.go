package match

import (
	"fmt"
	"reflect"
)

// TypeOf returns a new matcher that will match if the type names match
func TypeOf(value string) SupportedKindsMatcher {
	return &typeOf{value}
}

type typeOf struct {
	typeName string
}

// SupportedKinds returns all the kinds the type of matcher supports
func (m *typeOf) SupportedKinds() map[reflect.Kind]struct{} {
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

// Match returns true if actual type names match
// if the actual value is a pointer the type of what it points to will be used
func (m *typeOf) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	actualType := reflect.TypeOf(value)
	for actualType.Kind() == reflect.Ptr {
		actualType = actualType.Elem()
	}

	switch actualType.Kind() {
	case reflect.Slice:
		return "slice" == m.typeName || fmt.Sprintf("[]%v", actualType.Elem().Name()) == m.typeName
	case reflect.Array:
		return "array" == m.typeName || fmt.Sprintf("[%v]%v", actualType.Len(), actualType.Elem().Name()) == m.typeName
	default:
		return actualType.Name() == m.typeName
	}
}
