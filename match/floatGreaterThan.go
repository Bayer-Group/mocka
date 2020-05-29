package match

import (
	"reflect"
)

// FloatGreaterThan returns a new matcher that will match float's greater than the provided float
func FloatGreaterThan(value float64) SupportedKindsMatcher {
	return &floatGreaterThan{value}
}

type floatGreaterThan struct {
	value float64
}

// SupportedKinds returns all the kinds the float greater than matcher supports
func (m *floatGreaterThan) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Float32: {},
		reflect.Float64: {},
	}
}

// Match returns true if actual is an float greater than the provided float
func (m *floatGreaterThan) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.Float32:
		return value.(float32) > float32(m.value)
	case reflect.Float64:
		return value.(float64) > m.value
	default:
		return false
	}
}
