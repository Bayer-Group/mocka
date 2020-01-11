package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("anything", func() {
	Describe("Anything", func() {
		It("returns a anything struct", func() {
			actual := Anything()

			Expect(actual).To(BeAssignableToTypeOf(new(anything)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := Anything().SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
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
				}))
		})
	})

	DescribeTable("Match",
		func(input interface{}) {
			Expect(Anything().Match(input)).To(BeTrue())
		},
		Entry("returns true with nil", nil),
		Entry("returns true with number", 123),
		Entry("returns true with bool", true),
		Entry("returns true with struct", struct{}{}),
		Entry("returns true with chan", make(chan int, 1)),
		Entry("returns true with func", func() {}),
		Entry("returns true with pointer", &anything{}),
		Entry("returns true with slice", []string{"screams"}),
		Entry("returns true with array", [1]int{1}),
	)
})
