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
func (m *elementsContaining) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Slice: struct{}{},
		reflect.Array: struct{}{},
	}
}

// Match return true if the elements exist in the slice or array
func (m *elementsContaining) Match(value interface{}) bool {
	if value == nil {
		return false
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Slice, reflect.Array:
		allExist := true
		elementKind := v.Type().Elem().Kind()
		for _, e := range m.elements {
			if reflect.TypeOf(e).Kind() != elementKind {
				return false
			}

			found := false
			for i := 0; i < v.Len(); i++ {
				actual := v.Index(i)
				if reflect.DeepEqual(e, actual.Interface()) {
					found = true
					break
				}
			}

			allExist = allExist && found
		}

		return allExist
	default:
		return false
	}
}
