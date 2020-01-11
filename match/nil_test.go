package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("nil", func() {
	Describe("nil", func() {
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

	DescribeTable("Match true when",
		func(value interface{}) {
			Expect(Nil().Match(value)).To(BeTrue())
		},
		Entry("returns true with nil chan", (chan int)(nil)),
		Entry("returns true with nil func", (func())(nil)),
		Entry("returns true with nil interface", (SupportedKindsMatcher)((*nilMatcher)(nil))),
		Entry("returns true with nil map", (map[string]int)(nil)),
		Entry("returns true with nil pointer", (*anythingButNil)(nil)),
		Entry("returns true with nil slice", ([]string)(nil)),
	)

	DescribeTable("Match false when",
		func(value interface{}) {
			Expect(Nil().Match(value)).To(BeFalse())
		},
		Entry("returns false with non nil chan", make(chan int)),
		Entry("returns false with non nil func", func() {}),
		Entry("returns false with non nil interface", new(SupportedKindsMatcher)),
		Entry("returns false with non nil map", map[int]string{}),
		Entry("returns false with non nil pointer", &nilMatcher{}),
		Entry("returns false with non nil slice", []string{}),
	)
})
