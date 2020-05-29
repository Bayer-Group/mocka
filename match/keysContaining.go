package match

import (
	"reflect"
)

// KeysContaining returns a new matcher that will match the
// existence of map keys
func KeysContaining(keys ...interface{}) SupportedKindsMatcher {
	return &keysContaining{keys}
}

type keysContaining struct {
	keys []interface{}
}

// SupportedKinds returns all the kinds the keys containing matcher supports
func (keysContaining) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Map: {},
	}
}

// Match return true if the map contains the provided keys
func (m *keysContaining) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Map:
		keyKind := v.Type().Key().Kind()
		mapKeys := v.MapKeys()
		for _, k := range m.keys {
			if reflect.TypeOf(k).Kind() != keyKind {
				return false
			}

			var found bool
			for _, mk := range mapKeys {
				if reflect.DeepEqual(k, mk.Interface()) {
					found = true
					break
				}
			}

			if !found {
				return false
			}
		}

		return true
	default:
		return false
	}
}
