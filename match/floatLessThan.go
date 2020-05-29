package match

import (
	"reflect"
)

// FloatLessThan returns a new matcher that will match float's less than the provided float
func FloatLessThan(value float32) SupportedKindsMatcher {
	return &floatLessThan{value}
}

type floatLessThan struct {
	value float32
}

// SupportedKinds returns all the kinds the float less than matcher supports
func (m *floatLessThan) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Float32: {},
		reflect.Float64: {},
	}
}

// Match returns true if actual is an float less than the provided float
func (m *floatLessThan) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.Float32:
		return value.(float32) < m.value
	case reflect.Float64:
		return value.(float64) < float64(m.value)
	default:
		return false
	}
}
