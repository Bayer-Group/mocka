package match

import "reflect"

// Nil returns a new matcher that will only match nil
func Nil() SupportedKindsMatcher {
	return &nilMatcher{}
}

type nilMatcher struct {
}

// SupportedKinds returns all the kinds the nil matcher supports
func (m *nilMatcher) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Chan:      {},
		reflect.Func:      {},
		reflect.Interface: {},
		reflect.Map:       {},
		reflect.Ptr:       {},
		reflect.Slice:     {},
	}
}

// Match return true if the value is valid and nil
func (m *nilMatcher) Match(value interface{}) bool {
	v := reflect.ValueOf(value)
	return v.IsValid() && v.IsNil()
}
