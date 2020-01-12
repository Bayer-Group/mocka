package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("lengthOf", func() {
	Describe("LengthOth", func() {
		It("returns an lengthOf struct", func() {
			actual := LengthOf(2)

			Expect(actual).To(BeAssignableToTypeOf(new(lengthOf)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := LengthOf(3).SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.Array:  struct{}{},
					reflect.Map:    struct{}{},
					reflect.Slice:  struct{}{},
					reflect.String: struct{}{},
				}))
		})
	})

	DescribeTable("Match returns true",
		func(length int, actual interface{}) {
			Expect(LengthOf(length).Match(actual)).To(BeTrue())
		},
		Entry("when length matches for slice", 3, []int{1, 2, 3}),
		Entry("when length matches for array", 2, [2]string{"a", "b"}),
		Entry("when length matches for string", 5, "hello"),
		Entry("when length matches for map", 1, map[int]string{0: "a"}),
		Entry("with empty slice", 0, []int{}),
		Entry("with empty array", 0, [0]string{}),
		Entry("with empty string", 0, ""),
		Entry("with empty map", 0, map[int]string{}),
	)

	DescribeTable("Match returns false",
		func(length int, actual interface{}) {
			Expect(LengthOf(length).Match(actual)).To(BeFalse())
		},
		Entry("when actual is nil", 0, nil),
		Entry("when actual is not a valid kind", 1, 123),
		Entry("when length does not matches for slice", 4, []int{1, 2, 3}),
		Entry("when length does not matches for array", 1, [2]string{"a", "b"}),
		Entry("when length does not matches for string", 8, "hello"),
		Entry("when length does not matches for map", 2, map[int]string{0: "a"}),
	)
})
