package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("elementsContaining", func() {
	Describe("ElementsContaining", func() {
		It("returns an elementsContaining struct", func() {
			actual := ElementsContaining(2)

			Expect(actual).To(BeAssignableToTypeOf(new(elementsContaining)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := ElementsContaining(3).SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.Array: struct{}{},
					reflect.Slice: struct{}{},
				}))
		})
	})

	DescribeTable("Match returns true",
		func(elements []interface{}, actual interface{}) {
			Expect(ElementsContaining(elements...).Match(actual)).To(BeTrue())
		},
		Entry("when the elements and slice are empty", []interface{}{}, []int{}),
		Entry("when the elements and array are empty", []interface{}{}, [0]string{}),
		Entry("when all elements exist in slice", []interface{}{1, 2, 3}, []int{
			1, 2, 3, 4, 5,
		}),
		Entry("when all elements exist in array", []interface{}{"1", "2", "3"}, [5]string{
			"1", "2", "3", "4", "5",
		}),
	)

	DescribeTable("Match returns false",
		func(elements []interface{}, actual interface{}) {
			Expect(ElementsContaining(elements...).Match(actual)).To(BeFalse())
		},
		Entry("when actual is nil", []interface{}{1}, nil),
		Entry("when actual is not a valid kind", []interface{}{"1"}, 123),
		Entry("when provided elements are of the wrong kind", []interface{}{"a", "b"}, []int{1}),
		Entry("when all elements do not exist in the slice", []interface{}{1, 2, 3}, []int{
			1, 2, 4, 5,
		}),
		Entry("when all elements do not exist in the array", []interface{}{"1", "2", "3"}, [4]string{
			"1", "2", "4", "5",
		}),
	)
})
