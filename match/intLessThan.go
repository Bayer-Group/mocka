package match

import (
	"reflect"
)

// IntLessThan returns a new matcher that will match int's less than the provided int
func IntLessThan(value int) SupportedKindsMatcher {
	return &intLessThan{value}
}

type intLessThan struct {
	value int
}

// SupportedKinds returns all the kinds the int less than matcher supports
func (m *intLessThan) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Int:    {},
		reflect.Int8:   {},
		reflect.Int16:  {},
		reflect.Int32:  {},
		reflect.Int64:  {},
		reflect.Uint:   {},
		reflect.Uint8:  {},
		reflect.Uint16: {},
		reflect.Uint32: {},
		reflect.Uint64: {},
	}
}

// Match returns true if actual is an int less than the provided int
func (m *intLessThan) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.Int:
		return value.(int) < m.value
	case reflect.Int8:
		return value.(int8) < int8(m.value)
	case reflect.Int16:
		return value.(int16) < int16(m.value)
	case reflect.Int32:
		return value.(int32) < int32(m.value)
	case reflect.Int64:
		return value.(int64) < int64(m.value)
	case reflect.Uint:
		return value.(uint) < uint(m.value)
	case reflect.Uint8:
		return value.(uint8) < uint8(m.value)
	case reflect.Uint16:
		return value.(uint16) < uint16(m.value)
	case reflect.Uint32:
		return value.(uint32) < uint32(m.value)
	case reflect.Uint64:
		return value.(uint64) < uint64(m.value)
	default:
		return false
	}
}
