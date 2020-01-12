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
