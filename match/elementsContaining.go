package match

import (
	"reflect"
)

// ElementsContaining returns a new matcher that will match the
// existence of elements in a slice or array
func ElementsContaining(elements ...interface{}) SupportedKindsMatcher {
	return &elementsContaining{elements}
}

type elementsContaining struct {
	elements []interface{}
}

// SupportedKinds returns all the kinds the elements containing matcher supports
func (elementsContaining) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Slice: {},
		reflect.Array: {},
	}
}

// Match return true if the elements exist in the slice or array
func (m *elementsContaining) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Slice, reflect.Array:
		elementKind := v.Type().Elem().Kind()
		for _, e := range m.elements {
			if reflect.TypeOf(e).Kind() != elementKind {
				return false
			}

			var found bool
			for i := 0; !found && i < v.Len(); i++ {
				found = reflect.DeepEqual(e, v.Index(i).Interface())
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
