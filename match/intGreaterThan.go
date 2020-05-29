package match

import (
	"reflect"
)

// IntGreaterThan returns a new matcher that will match int's greater than the provided int
func IntGreaterThan(value int64) SupportedKindsMatcher {
	return &intGreaterThan{value}
}

type intGreaterThan struct {
	value int64
}

// SupportedKinds returns all the kinds the int greater than matcher supports
func (intGreaterThan) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Int:   {},
		reflect.Int8:  {},
		reflect.Int16: {},
		reflect.Int32: {},
		reflect.Int64: {},
	}
}

// Match returns true if actual is an int greater than the provided int
func (m *intGreaterThan) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.Int:
		return value.(int) > int(m.value)
	case reflect.Int8:
		return value.(int8) > int8(m.value)
	case reflect.Int16:
		return value.(int16) > int16(m.value)
	case reflect.Int32:
		return value.(int32) > int32(m.value)
	case reflect.Int64:
		return value.(int64) > m.value
	default:
		return false
	}
}
