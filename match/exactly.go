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
func (m *exactly) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Bool:          struct{}{},
		reflect.Int:           struct{}{},
		reflect.Int8:          struct{}{},
		reflect.Int16:         struct{}{},
		reflect.Int32:         struct{}{},
		reflect.Int64:         struct{}{},
		reflect.Uint:          struct{}{},
		reflect.Uint8:         struct{}{},
		reflect.Uint16:        struct{}{},
		reflect.Uint32:        struct{}{},
		reflect.Uint64:        struct{}{},
		reflect.Uintptr:       struct{}{},
		reflect.Float32:       struct{}{},
		reflect.Float64:       struct{}{},
		reflect.Complex64:     struct{}{},
		reflect.Complex128:    struct{}{},
		reflect.Array:         struct{}{},
		reflect.Chan:          struct{}{},
		reflect.Func:          struct{}{},
		reflect.Interface:     struct{}{},
		reflect.Map:           struct{}{},
		reflect.Ptr:           struct{}{},
		reflect.Slice:         struct{}{},
		reflect.String:        struct{}{},
		reflect.Struct:        struct{}{},
		reflect.UnsafePointer: struct{}{},
	}
}

// Match returns true when the values are equal using reflect.DeepEqual
func (m *exactly) Match(value interface{}) bool {
	return reflect.DeepEqual(m.value, value)
}
