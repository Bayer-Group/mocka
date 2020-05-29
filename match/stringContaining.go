package match

import (
	"reflect"
	"strings"
)

// StringContaining returns a new matcher that will match is the substring is found
func StringContaining(substring string) SupportedKindsMatcher {
	return &stringContaining{substring}
}

type stringContaining struct {
	substring string
}

// SupportedKinds returns all the kinds the string containing matcher supports
func (stringContaining) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.String: {},
	}
}

// Match return true if the provided substring is found
func (m *stringContaining) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.String:
		return strings.Contains(value.(string), m.substring)
	default:
		return false
	}
}
