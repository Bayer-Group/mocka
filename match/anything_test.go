package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("anything", func() {
	Describe("Anything", func() {
		It("returns an anything struct", func() {
			actual := Anything()

			Expect(actual).To(BeAssignableToTypeOf(new(anything)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := Anything().SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
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
				}))
		})
	})

	DescribeTable("Match returns true",
		func(input interface{}) {
			Expect(Anything().Match(input)).To(BeTrue())
		},
		Entry("with nil", nil),
		Entry("with number", 123),
		Entry("with bool", true),
		Entry("with struct", struct{}{}),
		Entry("with chan", make(chan int, 1)),
		Entry("with func", func() {}),
		Entry("with pointer", &anything{}),
		Entry("with slice", []string{"screams"}),
		Entry("with array", [1]int{1}),
	)
})
