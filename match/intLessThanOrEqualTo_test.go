package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("intLessThanOrEqualTo", func() {
	Describe("IntLessThanOrEqualTo", func() {
		It("returns an intLessThanOrEqualTo struct", func() {
			actual := IntLessThanOrEqualTo(10)

			Expect(actual).To(BeAssignableToTypeOf(new(intLessThanOrEqualTo)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := IntLessThanOrEqualTo(5).SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.Int:    {},
					reflect.Int8:   {},
					reflect.Int16:  {},
					reflect.Int32:  {},
					reflect.Int64:  {},
					reflect.Uint:   {},
					reflect.Uint8:  {},
					reflect.Uint16: {},
					reflect.Uint32: {},
					reflect.Uint64: {},
				}))
		})
	})

	DescribeTable("Match returns true",
		func(expected int, actual interface{}) {
			Expect(IntLessThanOrEqualTo(expected).Match(actual)).To(BeTrue())
		},
		Entry("with int", 10, int(5)),
		Entry("with int8", 18, int8(10)),
		Entry("with int16", 22, int16(15)),
		Entry("with int32", 40, int32(20)),
		Entry("with int64", 15, int64(8)),
		Entry("with uint", 10, uint(5)),
		Entry("with uint8", 18, uint8(10)),
		Entry("with uint16", 22, uint16(15)),
		Entry("with uint32", 40, uint32(20)),
		Entry("with uint64", 15, uint64(8)),
		Entry("when actual(int) is the same as the expected", 5, int(5)),
		Entry("when actual(int8) is the same as the expected", 10, int8(10)),
		Entry("when actual(int16) is the same as the expected", 15, int16(15)),
		Entry("when actual(int32) is the same as the expected", 20, int32(20)),
		Entry("when actual(int64) is the same as the expected", 8, int64(8)),
		Entry("when actual(uint) is the same as the expected", 5, uint(5)),
		Entry("when actual(uint8) is the same as the expected", 10, uint8(10)),
		Entry("when actual(uint16) is the same as the expected", 15, uint16(15)),
		Entry("when actual(uint32) is the same as the expected", 20, uint32(20)),
		Entry("when actual(uint64) is the same as the expected", 8, uint64(8)),
	)

	DescribeTable("Match returns false",
		func(expected int, actual interface{}) {
			Expect(IntLessThanOrEqualTo(expected).Match(actual)).To(BeFalse())
		},
		Entry("when actual is nil", 5, nil),
		Entry("when actual(int) is greater than expected", 1, int(5)),
		Entry("when actual(int8) is greater than expected", 8, int8(10)),
		Entry("when actual(int16) is greater than expected", 2, int16(15)),
		Entry("when actual(int32) is greater than expected", 4, int32(20)),
		Entry("when actual(int64) is greater than expected", 5, int64(8)),
		Entry("when actual(uint) is greater than expected", 1, uint(5)),
		Entry("when actual(uint8) is greater than expected", 8, uint8(10)),
		Entry("when actual(uint16) is greater than expected", 2, uint16(15)),
		Entry("when actual(uint32) is greater than expected", 4, uint32(20)),
		Entry("when actual(uint64) is greater than expected", 5, uint64(8)),
		Entry("when actual is not an int", 10, "10"),
	)
})
