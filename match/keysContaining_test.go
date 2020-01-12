package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("keysContaining", func() {
	Describe("KeysContaining", func() {
		It("returns an keysContaining struct", func() {
			actual := KeysContaining(2)

			Expect(actual).To(BeAssignableToTypeOf(new(keysContaining)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := KeysContaining(3).SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.Map: struct{}{},
				}))
		})
	})

	DescribeTable("Match returns true",
		func(keys []interface{}, actual interface{}) {
			Expect(KeysContaining(keys...).Match(actual)).To(BeTrue())
		},
		Entry("when the keys and map are empty", []interface{}{}, map[string]string{}),
		Entry("when all keys exist in the map", []interface{}{1, 2, 3}, map[int]string{
			1: "1",
			2: "2",
			3: "3",
			4: "4",
			5: "5",
		}),
	)

	DescribeTable("Match returns false",
		func(keys []interface{}, actual interface{}) {
			Expect(KeysContaining(keys...).Match(actual)).To(BeFalse())
		},
		Entry("when actual is nil", []interface{}{1}, nil),
		Entry("when actual is not a valid kind", []interface{}{"1"}, 123),
		Entry("when provided keys are of the wrong kind", []interface{}{"a", "b"}, map[int]string{}),
		Entry("when all keys do not exist", []interface{}{1, 2, 3}, map[int]string{
			1: "1",
			2: "2",
			4: "4",
			5: "5",
		}),
	)
})
