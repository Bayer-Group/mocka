package match

import (
	"reflect"
	"strings"
)

// StringPrefix returns a new matcher that will match the provided prefix
func StringPrefix(prefix string) SupportedKindsMatcher {
	return &stringPrefix{prefix}
}

type stringPrefix struct {
	prefix string
}

// SupportedKinds returns all the kinds the string prefix matcher supports
func (stringPrefix) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.String: {},
	}
}

// Match return true if the provided prefix is found
func (m *stringPrefix) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.String:
		return strings.HasPrefix(value.(string), m.prefix)
	default:
		return false
	}
}
