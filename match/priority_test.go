package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("priority", func() {
	DescribeTable("returns priority",
		func(matcher SupportedKindsMatcher, actual float64) {
			Expect(Priority(matcher)).To(Equal(actual))
		},
		Entry("priority for custom matchers", new(mockMatcher), float64(27)),
		Entry("priority for the exactly matcher", new(exactly), float64(25)),
		Entry("priority for the nilMatcher matcher", new(nilMatcher), float64(24)),
		Entry("priority for the floatGreaterThan matcher", new(floatGreaterThan), float64(23)),
		Entry("priority for the floatLessThan matcher", new(floatLessThan), float64(22)),
		Entry("priority for the floatGreaterThanOrEqualTo matcher", new(floatGreaterThanOrEqualTo), float64(21)),
		Entry("priority for the floatLessThanOrEqualTo matcher", new(floatLessThanOrEqualTo), float64(20)),
		Entry("priority for the intGreaterThan matcher", new(intGreaterThan), float64(19)),
		Entry("priority for the intLessThan matcher", new(intLessThan), float64(18)),
		Entry("priority for the intGreaterThanOrEqualTo matcher", new(intGreaterThanOrEqualTo), float64(17)),
		Entry("priority for the intLessThanOrEqualTo matcher", new(intLessThanOrEqualTo), float64(16)),
		Entry("priority for the uintGreaterThan matcher", new(uintGreaterThan), float64(15)),
		Entry("priority for the uintLessThan matcher", new(uintLessThan), float64(14)),
		Entry("priority for the uintGreaterThanOrEqualTo matcher", new(uintGreaterThanOrEqualTo), float64(13)),
		Entry("priority for the uintLessThanOrEqualTo matcher", new(uintLessThanOrEqualTo), float64(12)),
		Entry("priority for the stringPrefix matcher", new(stringPrefix), float64(11)),
		Entry("priority for the stringSuffix matcher", new(stringSuffix), float64(10)),
		Entry("priority for the stringContaining matcher", new(stringContaining), float64(9)),
		Entry("priority for the lengthOf matcher", new(lengthOf), float64(8)),
		Entry("priority for the empty matcher", new(empty), float64(7)),
		Entry("priority for the keysContaining matcher", new(keysContaining), float64(6)),
		Entry("priority for the elementsContaining matcher", new(elementsContaining), float64(5)),
		Entry("priority for the implementerOf matcher", new(implementerOf), float64(4)),
		Entry("priority for the convertibleTo matcher", new(convertibleTo), float64(3)),
		Entry("priority for the typeOf matcher", new(typeOf), float64(2)),
		Entry("priority for the anythingButNil matcher", new(anythingButNil), float64(1)),
		Entry("priority for the anything matcher", new(anything), float64(0)),
	)
})

type mockMatcher struct {
}

// SupportedKinds returns the supported kinds for the matcher
func (mockMatcher) SupportedKinds() map[reflect.Kind]struct{} {
	return nil
}

// Match return true is the match was successful; otherwise false
func (mockMatcher) Match(interface{}) bool {
	return true
}
