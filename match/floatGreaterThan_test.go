package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("floatGreaterThan", func() {
	Describe("FloatGreaterThan", func() {
		It("returns an floatGreaterThan struct", func() {
			actual := FloatGreaterThan(10)

			Expect(actual).To(BeAssignableToTypeOf(new(floatGreaterThan)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := FloatGreaterThan(5).SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.Float32: struct{}{},
					reflect.Float64: struct{}{},
				}))
		})
	})

	DescribeTable("Match returns true",
		func(expected float32, actual interface{}) {
			Expect(FloatGreaterThan(expected).Match(actual)).To(BeTrue())
		},
		Entry("with float32", float32(20), float32(40)),
		Entry("with float64", float32(8), float64(15)),
	)

	DescribeTable("Match returns false",
		func(expected float32, actual interface{}) {
			Expect(FloatGreaterThan(expected).Match(actual)).To(BeFalse())
		},
		Entry("when actual is nil", float32(5), nil),
		Entry("when actual(float32) is less than expected", float32(20), float32(4)),
		Entry("when actual(float64) is less than expected", float32(8), float64(5)),
		Entry("when actual(float32) is the same as the expected", float32(20), float32(20)),
		Entry("when actual(float64) is the same as the expected", float32(8), float64(8)),
		Entry("when actual is not an int", float32(10), "10"),
	)
})
