package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("intGreaterThan", func() {
	Describe("IntGreaterThan", func() {
		It("returns an intGreaterThan struct", func() {
			actual := IntGreaterThan(10)

			Expect(actual).To(BeAssignableToTypeOf(new(intGreaterThan)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := IntGreaterThan(5).SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.Int:   struct{}{},
					reflect.Int8:  struct{}{},
					reflect.Int16: struct{}{},
					reflect.Int32: struct{}{},
					reflect.Int64: struct{}{},
				}))
		})
	})

	DescribeTable("Match returns true",
		func(expected int, actual interface{}) {
			Expect(IntGreaterThan(expected).Match(actual)).To(BeTrue())
		},
		Entry("with int", 5, int(10)),
		Entry("with int8", 10, int8(18)),
		Entry("with int16", 15, int16(22)),
		Entry("with int32", 20, int32(40)),
		Entry("with int64", 8, int64(15)),
	)

	DescribeTable("Match returns false",
		func(expected int, actual interface{}) {
			Expect(IntGreaterThan(expected).Match(actual)).To(BeFalse())
		},
		Entry("when actual is nil", 5, nil),
		Entry("when int is less than expected", 5, int(1)),
		Entry("when int8 is less than expected", 10, int8(8)),
		Entry("when int16 is less than expected", 15, int16(2)),
		Entry("when int32 is less than expected", 20, int32(4)),
		Entry("when int64 is less than expected", 8, int64(5)),
		Entry("when actual is not an int", 10, "10"),
	)
})
