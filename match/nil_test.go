package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("nil", func() {
	Describe("Nil", func() {
		It("returns an exactly struct", func() {
			actual := Nil()

			Expect(actual).To(BeAssignableToTypeOf(new(nilMatcher)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := Nil().SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.Chan:      struct{}{},
					reflect.Func:      struct{}{},
					reflect.Interface: struct{}{},
					reflect.Map:       struct{}{},
					reflect.Ptr:       struct{}{},
					reflect.Slice:     struct{}{},
				}))
		})
	})

	DescribeTable("Match returns true",
		func(value interface{}) {
			Expect(Nil().Match(value)).To(BeTrue())
		},
		Entry("with nil chan", (chan int)(nil)),
		Entry("with nil func", (func())(nil)),
		Entry("with nil interface", (SupportedKindsMatcher)((*nilMatcher)(nil))),
		Entry("with nil map", (map[string]int)(nil)),
		Entry("with nil pointer", (*anythingButNil)(nil)),
		Entry("with nil slice", ([]string)(nil)),
	)

	DescribeTable("Match returns false",
		func(value interface{}) {
			Expect(Nil().Match(value)).To(BeFalse())
		},
		Entry("with non nil chan", make(chan int)),
		Entry("with non nil func", func() {}),
		Entry("with non nil interface", new(SupportedKindsMatcher)),
		Entry("with non nil map", map[int]string{}),
		Entry("with non nil pointer", &nilMatcher{}),
		Entry("with non nil slice", []string{}),
	)
})
