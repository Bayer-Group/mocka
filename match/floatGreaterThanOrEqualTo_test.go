package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("floatGreaterThanOrEqualTo", func() {
	Describe("FloatGreaterThanOrEqualTo", func() {
		It("returns an floatGreaterThanOrEqualTo struct", func() {
			actual := FloatGreaterThanOrEqualTo(10)

			Expect(actual).To(BeAssignableToTypeOf(new(floatGreaterThanOrEqualTo)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := FloatGreaterThanOrEqualTo(5).SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.Float32: {},
					reflect.Float64: {},
				}))
		})
	})

	DescribeTable("Match returns true",
		func(expected float64, actual interface{}) {
			Expect(FloatGreaterThanOrEqualTo(expected).Match(actual)).To(BeTrue())
		},
		Entry("with float32", float64(20), float32(40)),
		Entry("with float64", float64(8), float64(15)),
		Entry("when actual(float32) is the same as the expected", float64(20), float32(20)),
		Entry("when actual(float64) is the same as the expected", float64(8), float64(8)),
	)

	DescribeTable("Match returns false",
		func(expected float64, actual interface{}) {
			Expect(FloatGreaterThanOrEqualTo(expected).Match(actual)).To(BeFalse())
		},
		Entry("when actual is nil", float64(5), nil),
		Entry("when actual(float32) is less than expected", float64(20), float32(4)),
		Entry("when actual(float64) is less than expected", float64(8), float64(5)),
		Entry("when actual is not an int", float64(10), "10"),
	)
})
