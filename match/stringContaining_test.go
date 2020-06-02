package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("stringContaining", func() {
	Describe("StringContaining", func() {
		It("returns a stringContaining struct", func() {
			actual := StringContaining("")

			Expect(actual).To(BeAssignableToTypeOf(new(stringContaining)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := StringContaining("").SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.String: {},
				}))
		})
	})

	DescribeTable("Match returns true",
		func(expected string, actual interface{}) {
			Expect(StringContaining(expected).Match(actual)).To(BeTrue())
		},
		Entry("when the string contains the substring", "sub", "I have a substring"),
	)

	DescribeTable("Match returns false",
		func(expected string, actual interface{}) {
			Expect(StringContaining(expected).Match(actual)).To(BeFalse())
		},
		Entry("when actual is nil", "hi", nil),
		Entry("when actual is not a string", "hi", 12),
		Entry("when the substring does not exist in actual", "hello", "screams"),
	)
})
