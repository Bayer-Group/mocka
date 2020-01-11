package match

import "reflect"

// Anything returns a new matcher that will match any value
func Anything() SupportedKindsMatcher {
	return &anything{}
}

type anything struct {
}

// SupportedKinds returns all the kinds the anything matcher supports
func (m *anything) SupportedKinds() map[reflect.Kind]struct{} {
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

// Match always returns true
func (m *anything) Match(_ interface{}) bool {
	return true
}
