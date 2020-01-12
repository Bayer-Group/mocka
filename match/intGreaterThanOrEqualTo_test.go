package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("intGreaterThanOrEqualTo", func() {
	Describe("IntGreaterThanOrEqualTo", func() {
		It("returns an intGreaterThanOrEqualTo struct", func() {
			actual := IntGreaterThanOrEqualTo(10)

			Expect(actual).To(BeAssignableToTypeOf(new(intGreaterThanOrEqualTo)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := IntGreaterThanOrEqualTo(5).SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.Int:    struct{}{},
					reflect.Int8:   struct{}{},
					reflect.Int16:  struct{}{},
					reflect.Int32:  struct{}{},
					reflect.Int64:  struct{}{},
					reflect.Uint:   struct{}{},
					reflect.Uint8:  struct{}{},
					reflect.Uint16: struct{}{},
					reflect.Uint32: struct{}{},
					reflect.Uint64: struct{}{},
				}))
		})
	})

	DescribeTable("Match returns true",
		func(expected int, actual interface{}) {
			Expect(IntGreaterThanOrEqualTo(expected).Match(actual)).To(BeTrue())
		},
		Entry("with int", 5, int(10)),
		Entry("with int8", 10, int8(18)),
		Entry("with int16", 15, int16(22)),
		Entry("with int32", 20, int32(40)),
		Entry("with int64", 8, int64(15)),
		Entry("with uint", 5, uint(10)),
		Entry("with uint8", 10, uint8(18)),
		Entry("with uint16", 15, uint16(22)),
		Entry("with uint32", 20, uint32(40)),
		Entry("with uint64", 8, uint64(15)),
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
			Expect(IntGreaterThanOrEqualTo(expected).Match(actual)).To(BeFalse())
		},
		Entry("when actual is nil", 5, nil),
		Entry("when actual(int) is less than expected", 5, int(1)),
		Entry("when actual(int8) is less than expected", 10, int8(8)),
		Entry("when actual(int16) is less than expected", 15, int16(2)),
		Entry("when actual(int32) is less than expected", 20, int32(4)),
		Entry("when actual(int64) is less than expected", 8, int64(5)),
		Entry("when actual(uint) is less than expected", 5, uint(1)),
		Entry("when actual(uint8) is less than expected", 10, uint8(8)),
		Entry("when actual(uint16) is less than expected", 15, uint16(2)),
		Entry("when actual(uint32) is less than expected", 20, uint32(4)),
		Entry("when actual(uint64) is less than expected", 8, uint64(5)),
		Entry("when actual is not an int", 10, "10"),
	)
})
