package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("uintGreaterThanOrEqualTo", func() {
	Describe("UintGreaterThanOrEqualTo", func() {
		It("returns an uintGreaterThanOrEqualTo struct", func() {
			actual := UintGreaterThanOrEqualTo(uint64(10))

			Expect(actual).To(BeAssignableToTypeOf(new(uintGreaterThanOrEqualTo)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := UintGreaterThanOrEqualTo(uint64(5)).SupportedKinds()

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
			Expect(UintGreaterThanOrEqualTo(expected).Match(actual)).To(BeTrue())
		},
		Entry("with uint", uint64(5), uint(10)),
		Entry("with uint8", uint64(10), uint8(18)),
		Entry("with uint16", uint64(15), uint16(22)),
		Entry("with uint32", uint64(20), uint32(40)),
		Entry("with uint64", uint64(8), uint64(15)),
		Entry("when actual(uint) is the same as the expected", uint64(5), uint(5)),
		Entry("when actual(uint8) is the same as the expected", uint64(10), uint8(10)),
		Entry("when actual(uint16) is the same as the expected", uint64(15), uint16(15)),
		Entry("when actual(uint32) is the same as the expected", uint64(20), uint32(20)),
		Entry("when actual(uint64) is the same as the expected", uint64(8), uint64(8)),
	)

	DescribeTable("Match returns false",
		func(expected uint64, actual interface{}) {
			Expect(UintGreaterThanOrEqualTo(expected).Match(actual)).To(BeFalse())
		},
		Entry("when actual(uint) is less than expected", uint64(5), uint(1)),
		Entry("when actual(uint8) is less than expected", uint64(10), uint8(8)),
		Entry("when actual(uint16) is less than expected", uint64(15), uint16(2)),
		Entry("when actual(uint32) is less than expected", uint64(20), uint32(4)),
		Entry("when actual(uint64) is less than expected", uint64(8), uint64(5)),
		Entry("when actual is not an int", uint64(10), "10"),
	)
})
