package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("anythingButNil", func() {
	Describe("AnythingButNil", func() {
		It("returns an anythingButNil struct", func() {
			actual := AnythingButNil()

			Expect(actual).To(BeAssignableToTypeOf(new(anythingButNil)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := AnythingButNil().SupportedKinds()

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
			Expect(AnythingButNil().Match(value)).To(BeTrue())
		},
		Entry("with non nil chan", make(chan int)),
		Entry("with non nil func", func() {}),
		Entry("with non nil interface", new(SupportedKindsMatcher)),
		Entry("with non nil map", map[int]string{}),
		Entry("with non nil pointer", &anythingButNil{}),
		Entry("with non nil slice", []string{}),
	)

	DescribeTable("Match returns false",
		func(value interface{}) {
			Expect(AnythingButNil().Match(value)).To(BeFalse())
		},
		Entry("with nil chan", (chan int)(nil)),
		Entry("with nil func", (func())(nil)),
		Entry("with nil interface", (SupportedKindsMatcher)((*anythingButNil)(nil))),
		Entry("with nil map", (map[string]int)(nil)),
		Entry("with nil pointer", (*anythingButNil)(nil)),
		Entry("with nil slice", ([]string)(nil)),
	)
})
