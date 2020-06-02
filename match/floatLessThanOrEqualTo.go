package match

import (
	"reflect"
)

// FloatLessThanOrEqualTo returns a new matcher that will match float's less than or equal to the provided float
func FloatLessThanOrEqualTo(value float64) SupportedKindsMatcher {
	return &floatLessThanOrEqualTo{value}
}

type floatLessThanOrEqualTo struct {
	value float64
}

// SupportedKinds returns all the kinds the float less than or equal to matcher supports
func (floatLessThanOrEqualTo) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Float32: {},
		reflect.Float64: {},
	}
}

// Match returns true if actual is an float less than or equal to the provided float
func (m *floatLessThanOrEqualTo) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.Float32:
		return value.(float32) <= float32(m.value)
	case reflect.Float64:
		return value.(float64) <= m.value
	default:
		return false
	}
}
