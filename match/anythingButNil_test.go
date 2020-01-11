package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("anythingButNil", func() {
	Describe("AnythingButNil", func() {
		It("returns an exactly struct", func() {
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

	DescribeTable("Match true when",
		func(value interface{}) {
			Expect(AnythingButNil().Match(value)).To(BeTrue())
		},
		Entry("returns true with non nil chan", make(chan int)),
		Entry("returns true with non nil func", func() {}),
		Entry("returns true with non nil interface", new(SupportedKindsMatcher)),
		Entry("returns true with non nil map", map[int]string{}),
		Entry("returns true with non nil pointer", &anythingButNil{}),
		Entry("returns true with non nil slice", []string{}),
	)

	DescribeTable("Match true when",
		func(value interface{}) {
			Expect(AnythingButNil().Match(value)).To(BeFalse())
		},
		Entry("returns false with nil chan", (chan int)(nil)),
		Entry("returns false with nil func", (func())(nil)),
		Entry("returns false with nil interface", (SupportedKindsMatcher)(nil)),
		Entry("returns false with nil map", (map[string]int)(nil)),
		Entry("returns false with nil pointer", (*anythingButNil)(nil)),
		Entry("returns false with nil slice", ([]string)(nil)),
	)
})
