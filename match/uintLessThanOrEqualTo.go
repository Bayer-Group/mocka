package match

import (
	"reflect"
)

// UintLessThanOrEqualTo returns a new matcher that will match uint's less than or equal to the provided uint
func UintLessThanOrEqualTo(value uint64) SupportedKindsMatcher {
	return &uintLessThanOrEqualTo{value}
}

type uintLessThanOrEqualTo struct {
	value uint64
}

// SupportedKinds returns all the kinds the int less than or equal to matcher supports
func (uintLessThanOrEqualTo) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Uint:   {},
		reflect.Uint8:  {},
		reflect.Uint16: {},
		reflect.Uint32: {},
		reflect.Uint64: {},
	}
}

// Match returns true if actual is an int less than or equal to the provided int
func (m *uintLessThanOrEqualTo) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.Uint:
		return value.(uint) <= uint(m.value)
	case reflect.Uint8:
		return value.(uint8) <= uint8(m.value)
	case reflect.Uint16:
		return value.(uint16) <= uint16(m.value)
	case reflect.Uint32:
		return value.(uint32) <= uint32(m.value)
	case reflect.Uint64:
		return value.(uint64) <= uint64(m.value)
	default:
		return false
	}
}
