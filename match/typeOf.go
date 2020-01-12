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
