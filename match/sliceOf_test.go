package match

import (
	"errors"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("sliceOf", func() {
	var matcher SupportedKindsMatcher

	BeforeEach(func() {
		matcher = SliceOf(Anything(), Nil(), Exactly("A"))
	})

	Describe("SliceOf", func() {
		It("returns a sliceOf struct", func() {
			Expect(matcher).To(BeAssignableToTypeOf(new(sliceOf)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			Expect(matcher.SupportedKinds()).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.Array: {},
					reflect.Slice: {},
				}))
		})
	})

	Describe("Match", func() {
		It("returns true is all matchers are truthy", func() {
			Expect(matcher.Match([]interface{}{1, nil, "A"})).To(BeTrue())
		})

		It("return false if the length of arguments do not match the length of matchers", func() {
			Expect(matcher.Match([]interface{}{1})).To(BeFalse())
		})

		It("returns false if one of the matchers is not truthy", func() {
			Expect(matcher.Match([]interface{}{1, errors.New("a"), "A"})).To(BeFalse())
		})
	})
})
