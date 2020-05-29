package match

import (
	"reflect"
)

// ImplementerOf returns a new matcher that will match a value if it implements the provided interface
func ImplementerOf(value interface{}) SupportedKindsMatcher {
	return &implementerOf{value}
}

type implementerOf struct {
	value interface{}
}

// SupportedKinds returns all the kinds the implementer of matcher supports
func (m *implementerOf) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Ptr: {},
	}
}

// Match returns true if the value implements the provided interface; otherwise false.
// If provided value is nil or is not an interface the matcher will always return false
func (m *implementerOf) Match(value interface{}) bool {
	expectedValue := reflect.ValueOf(m.value)
	actualValue := reflect.ValueOf(value)

	if !expectedValue.IsValid() || !actualValue.IsValid() {
		return false
	}

	if expectedValue.IsNil() || actualValue.IsNil() {
		return false
	}

	expectedType := expectedValue.Type()
	actualType := actualValue.Type()

	if expectedType.Kind() != reflect.Ptr && expectedType.Elem().Kind() != reflect.Interface {
		return false
	}

	if actualType.Kind() != reflect.Ptr && actualType.Elem().Kind() != reflect.Struct {
		return false
	}

	return actualType.Implements(expectedType.Elem())
}
