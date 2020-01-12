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
