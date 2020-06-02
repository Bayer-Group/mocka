package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("convertibleTo", func() {
	Describe("ConvertibleTo", func() {
		It("returns a convertibleTo struct", func() {
			actual := ConvertibleTo(new(int))

			Expect(actual).To(BeAssignableToTypeOf(new(convertibleTo)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := ConvertibleTo(new(int)).SupportedKinds()

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
		func(i interface{}, value interface{}) {
			Expect(ConvertibleTo(i).Match(value)).To(BeTrue())
		},
		Entry("with numbers", new(int), int8(5)),
		Entry("with strings", new([]int32), "hello"),
		Entry("with slices", new(string), int32(5)),
		Entry("with interfaces", new(SupportedKindsMatcher), &convertibleTo{}),
	)

	DescribeTable("Match returns false",
		func(i interface{}, value interface{}) {
			Expect(ConvertibleTo(i).Match(value)).To(BeFalse())
		},
		Entry("when expected is nil", nil, &convertibleTo{}),
		Entry("when actual is nil", new(SupportedKindsMatcher), nil),
		Entry("when expected is not an pointer", *new(int), 10),
		Entry("when actual cannot be converted to expected", new(int), []string{}),
	)
})
