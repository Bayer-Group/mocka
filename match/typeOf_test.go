package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("typeOf", func() {
	Describe("TypeOf", func() {
		It("returns an typeOf struct", func() {
			actual := TypeOf("int")

			Expect(actual).To(BeAssignableToTypeOf(new(typeOf)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := TypeOf("bool").SupportedKinds()

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

	DescribeTable("Match returns true when type names matches",
		func(typeName string, value interface{}) {
			Expect(TypeOf(typeName).Match(value)).To(BeTrue())
		},
		Entry("with numbers", "int", 5),
		Entry("with strings", "string", "hello"),
		Entry("with typed slices", "[]int", []int{1, 2, 3}),
		Entry("with typed arrays", "[3]int", [3]int{1, 2, 3}),
		Entry("with slices", "slice", []int{1, 2, 3}),
		Entry("with arrays", "array", [3]int{1, 2, 3}),
		Entry("with interfaces", "convertibleTo", &convertibleTo{}),
	)

	DescribeTable("Match returns false",
		func(typeName string, value interface{}) {
			Expect(TypeOf(typeName).Match(value)).To(BeFalse())
		},
		Entry("when actual is nil", "int", nil),
		Entry("when the type names do not match", "int", "i am an int"),
	)
})
