package match

import (
	"reflect"
)

// SliceOf returns a new matcher of matchers for slices and arrays
func SliceOf(matchers ...SupportedKindsMatcher) SupportedKindsMatcher {
	return &sliceOf{matchers}
}

type sliceOf struct {
	matchers []SupportedKindsMatcher
}

// SupportedKinds returns the kinds the sliceOf matcher supports
func (sliceOf) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Slice: {},
		reflect.Array: {},
	}
}

// Match returns true when all value elements match their respective matchers
func (m *sliceOf) Match(value interface{}) bool {
	slice := reflect.ValueOf(value)
	if _, exists := m.SupportedKinds()[slice.Kind()]; !exists {
		return false
	}

	if len(m.matchers) != slice.Len() {
		return false
	}

	for i, matcher := range m.matchers {
		indexValue := slice.Index(i)
		if _, exists := matcher.SupportedKinds()[indexValue.Kind()]; !exists {
			return false
		}

		var v interface{}
		if indexValue.IsValid() {
			v = indexValue.Interface()
		} else {
			v = nil
		}

		if !matcher.Match(v) {
			return false
		}
	}

	return true
}
