package match

import "reflect"

// Exactly returns a new matcher for matching exact values with reflect.DeepEqual
func Exactly(value interface{}) SupportedKindsMatcher {
	return &exactly{value}
}

type exactly struct {
	value interface{}
}

// SupportedKinds returns the kinds the exactly matcher supports
func (exactly) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Bool:          {},
		reflect.Int:           {},
		reflect.Int8:          {},
		reflect.Int16:         {},
		reflect.Int32:         {},
		reflect.Int64:         {},
		reflect.Uint:          {},
		reflect.Uint8:         {},
		reflect.Uint16:        {},
		reflect.Uint32:        {},
		reflect.Uint64:        {},
		reflect.Uintptr:       {},
		reflect.Float32:       {},
		reflect.Float64:       {},
		reflect.Complex64:     {},
		reflect.Complex128:    {},
		reflect.Array:         {},
		reflect.Chan:          {},
		reflect.Func:          {},
		reflect.Interface:     {},
		reflect.Map:           {},
		reflect.Ptr:           {},
		reflect.Slice:         {},
		reflect.String:        {},
		reflect.Struct:        {},
		reflect.UnsafePointer: {},
	}
}

// Match returns true when the values are equal using reflect.DeepEqual
func (m *exactly) Match(value interface{}) bool {
	return reflect.DeepEqual(m.value, value)
}
