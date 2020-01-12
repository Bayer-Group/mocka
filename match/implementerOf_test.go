package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("implementerOf", func() {
	Describe("ImplementerOf", func() {
		It("returns an exactly struct", func() {
			actual := ImplementerOf(new(SupportedKindsMatcher))

			Expect(actual).To(BeAssignableToTypeOf(new(implementerOf)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := ImplementerOf(new(SupportedKindsMatcher)).SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.Ptr: struct{}{},
				}))
		})
	})

	Describe("Match", func() {
		It("returns true when the struct implements the provided interface", func() {
			Expect(ImplementerOf(new(SupportedKindsMatcher)).Match(&implementerOf{})).To(BeTrue())
		})
	})

	DescribeTable("Match returns false",
		func(i interface{}, value interface{}) {
			Expect(ImplementerOf(i).Match(value)).To(BeFalse())
		},
		Entry("when interface is not valid", (SupportedKindsMatcher)(nil), &implementerOf{}),
		Entry("when value is not valid", new(SupportedKindsMatcher), (SupportedKindsMatcher)(nil)),
		Entry("when interface is nil", (SupportedKindsMatcher)((*implementerOf)(nil)), &implementerOf{}),
		Entry("when value is nil", new(SupportedKindsMatcher), (*implementerOf)(nil)),
		Entry("when interface is not an interface", make(chan int), &implementerOf{}),
		Entry("when value is not an interface", new(SupportedKindsMatcher), make(chan int)),
		Entry("when value does not implement the interface", new(SupportedKindsMatcher), &struct{}{}),
	)
})
