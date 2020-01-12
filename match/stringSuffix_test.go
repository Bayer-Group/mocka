package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("stringSuffix", func() {
	Describe("StringSuffix", func() {
		It("returns a stringSuffix struct", func() {
			actual := StringSuffix("")

			Expect(actual).To(BeAssignableToTypeOf(new(stringSuffix)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := StringSuffix("").SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.String: struct{}{},
				}))
		})
	})

	DescribeTable("Match returns true",
		func(expected string, actual interface{}) {
			Expect(StringSuffix(expected).Match(actual)).To(BeTrue())
		},
		Entry("when the string suffix is found", "a suffix", "I have a suffix"),
	)

	DescribeTable("Match returns false",
		func(expected string, actual interface{}) {
			Expect(StringSuffix(expected).Match(actual)).To(BeFalse())
		},
		Entry("when actual is nil", "hi", nil),
		Entry("when actual is not a string", "hi", 12),
		Entry("when the suffix does not exist in actual", "hello", "screams"),
	)
})
