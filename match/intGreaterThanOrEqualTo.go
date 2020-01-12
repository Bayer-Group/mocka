package match

import (
	"reflect"
)

// IntGreaterThanOrEqualTo returns a new matcher that will match int's greater than or equal to the provided int
func IntGreaterThanOrEqualTo(value int) SupportedKindsMatcher {
	return &intGreaterThanOrEqualTo{value}
}

type intGreaterThanOrEqualTo struct {
	value int
}

// SupportedKinds returns all the kinds the int greater than or equal to matcher supports
func (m *intGreaterThanOrEqualTo) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Int:   struct{}{},
		reflect.Int8:  struct{}{},
		reflect.Int16: struct{}{},
		reflect.Int32: struct{}{},
		reflect.Int64: struct{}{},
	}
}

// Match returns true if actual is an int greater than or equal to the provided int
func (m *intGreaterThanOrEqualTo) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.Int:
		return value.(int) >= m.value
	case reflect.Int8:
		return value.(int8) >= int8(m.value)
	case reflect.Int16:
		return value.(int16) >= int16(m.value)
	case reflect.Int32:
		return value.(int32) >= int32(m.value)
	case reflect.Int64:
		return value.(int64) >= int64(m.value)
	default:
		return false
	}
}
