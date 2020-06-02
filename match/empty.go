package match

import "reflect"

// Empty returns a new matcher that will match empty
// strings, slices, arrays, and maps
func Empty() SupportedKindsMatcher {
	return &empty{}
}

type empty struct {
}

// SupportedKinds returns all the kinds the empty matcher supports
func (empty) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Array:  {},
		reflect.Map:    {},
		reflect.Slice:  {},
		reflect.String: {},
	}
}

// Match return true if actual is an empty value; otherwise false
func (empty) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Array, reflect.Slice:
		return v.Len() == 0
	case reflect.Map:
		return len(v.MapKeys()) == 0
	case reflect.String:
		return len(value.(string)) == 0
	default:
		return false
	}
}
