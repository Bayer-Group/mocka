package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("stringPrefix", func() {
	Describe("StringPrefix", func() {
		It("returns a stringPrefix struct", func() {
			actual := StringPrefix("")

			Expect(actual).To(BeAssignableToTypeOf(new(stringPrefix)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := StringPrefix("").SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.String:      struct{}{},
				}))
		})
	})

	DescribeTable("Match returns true",
		func(expected string, actual interface{}) {
			Expect(StringPrefix(expected).Match(actual)).To(BeTrue())
		},
		Entry("when the prefix is found", "I am", "I am a string"),
	)

	DescribeTable("Match returns false",
		func(expected string, actual interface{}) {
			Expect(StringPrefix(expected).Match(actual)).To(BeFalse())
		},
		Entry("when actual is nil", "hi", nil),
		Entry("when actual is not a string", "hi", 12),
		Entry("when the prefix does not exist in actual", "hello", "screams"),
	)
})
