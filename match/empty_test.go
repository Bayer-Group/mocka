package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("empty", func() {
	Describe("Empty", func() {
		It("returns an empty struct", func() {
			actual := Empty()

			Expect(actual).To(BeAssignableToTypeOf(new(empty)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := Empty().SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.Array:  {},
					reflect.Map:    {},
					reflect.Slice:  {},
					reflect.String: {},
				}))
		})
	})

	DescribeTable("Match returns true",
		func(actual interface{}) {
			Expect(Empty().Match(actual)).To(BeTrue())
		},
		Entry("with empty slice", []int{}),
		Entry("with empty array", [0]string{}),
		Entry("with empty string", ""),
		Entry("with empty map", map[int]string{}),
	)

	DescribeTable("Match returns false",
		func(actual interface{}) {
			Expect(Empty().Match(actual)).To(BeFalse())
		},
		Entry("when actual is nil", nil),
		Entry("when actual is not a valid kind", 123),
		Entry("when length != 0 for slice", []int{1, 2, 3}),
		Entry("when length != 0 for array", [1]string{"a"}),
		Entry("when length != 0 for string", "hello"),
		Entry("when length != 0 for map", map[int]string{0: "a"}),
	)
})
