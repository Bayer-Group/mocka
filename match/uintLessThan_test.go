package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("uintLessThan", func() {
	Describe("UintLessThan", func() {
		It("returns an uintLessThan struct", func() {
			actual := UintLessThan(uint64(10))

			Expect(actual).To(BeAssignableToTypeOf(new(uintLessThan)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := UintLessThan(uint64(5)).SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.Uint:   {},
					reflect.Uint8:  {},
					reflect.Uint16: {},
					reflect.Uint32: {},
					reflect.Uint64: {},
				}))
		})
	})

	DescribeTable("Match returns true",
		func(expected uint64, actual interface{}) {
			Expect(UintLessThan(expected).Match(actual)).To(BeTrue())
		},
		Entry("with uint", uint64(10), uint(5)),
		Entry("with uint8", uint64(18), uint8(10)),
		Entry("with uint16", uint64(22), uint16(15)),
		Entry("with uint32", uint64(40), uint32(20)),
		Entry("with uint64", uint64(15), uint64(8)),
	)

	DescribeTable("Match returns false",
		func(expected uint64, actual interface{}) {
			Expect(UintLessThan(expected).Match(actual)).To(BeFalse())
		},
		Entry("when actual is nil", uint64(5), nil),
		Entry("when actual(uint) is greater than expected", uint64(1), uint(5)),
		Entry("when actual(uint8) is greater than expected", uint64(8), uint8(10)),
		Entry("when actual(uint16) is greater than expected", uint64(2), uint16(15)),
		Entry("when actual(uint32) is greater than expected", uint64(4), uint32(20)),
		Entry("when actual(uint64) is greater than expected", uint64(5), uint64(8)),
		Entry("when actual(uint) is the same as the expected", uint64(5), uint(5)),
		Entry("when actual(uint8) is the same as the expected", uint64(10), uint8(10)),
		Entry("when actual(uint16) is the same as the expected", uint64(15), uint16(15)),
		Entry("when actual(uint32) is the same as the expected", uint64(20), uint32(20)),
		Entry("when actual(uint64) is the same as the expected", uint64(8), uint64(8)),
		Entry("when actual is not an int", uint64(10), "10"),
	)
})
