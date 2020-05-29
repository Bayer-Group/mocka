package match

import "reflect"

// AnythingButNil returns a new matcher that will match any value except nil
func AnythingButNil() SupportedKindsMatcher {
	return &anythingButNil{}
}

type anythingButNil struct {
}

// SupportedKinds returns all the kinds the anything but nil matcher supports
func (m *anythingButNil) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Chan:      {},
		reflect.Func:      {},
		reflect.Interface: {},
		reflect.Map:       {},
		reflect.Ptr:       {},
		reflect.Slice:     {},
	}
}

// Match return true if the value is valid and not nil
func (m *anythingButNil) Match(value interface{}) bool {
	v := reflect.ValueOf(value)
	return v.IsValid() && !v.IsNil()
}
