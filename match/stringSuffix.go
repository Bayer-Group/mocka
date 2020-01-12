package match

import (
	"reflect"
	"strings"
)

// StringSuffix returns a new matcher that will match the provided suffix
func StringSuffix(suffix string) SupportedKindsMatcher {
	return &stringSuffix{suffix}
}

type stringSuffix struct {
	suffix string
}

// SupportedKinds returns all the kinds the string suffix matcher supports
func (m *stringSuffix) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.String: struct{}{},
	}
}

// Match return true if the provided suffix is found
func (m *stringSuffix) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.String:
		return strings.HasSuffix(value.(string), m.suffix)
	default:
		return false
	}
}
