package match

import "reflect"

// Priority returns the matchers priority to be compared against
func Priority(m SupportedKindsMatcher) float64 {
	if p, exists := priorities[reflect.TypeOf(m)]; exists {
		return p
	}

	return float64(len(priorities)) + 1
}

// priorities defines the priority ranking for custom matchers
var priorities = map[reflect.Type]float64{
	// exact value matchers
	reflect.TypeOf(new(exactly)):    25,
	reflect.TypeOf(new(nilMatcher)): 24,

	// numeric matchers
	reflect.TypeOf(new(floatGreaterThan)):          23,
	reflect.TypeOf(new(floatLessThan)):             22,
	reflect.TypeOf(new(floatGreaterThanOrEqualTo)): 21,
	reflect.TypeOf(new(floatLessThanOrEqualTo)):    20,

	reflect.TypeOf(new(intGreaterThan)):          19,
	reflect.TypeOf(new(intLessThan)):             18,
	reflect.TypeOf(new(intGreaterThanOrEqualTo)): 17,
	reflect.TypeOf(new(intLessThanOrEqualTo)):    16,

	reflect.TypeOf(new(uintGreaterThan)):          15,
	reflect.TypeOf(new(uintLessThan)):             14,
	reflect.TypeOf(new(uintGreaterThanOrEqualTo)): 13,
	reflect.TypeOf(new(uintLessThanOrEqualTo)):    12,

	// string matchers
	reflect.TypeOf(new(stringPrefix)):     11,
	reflect.TypeOf(new(stringSuffix)):     10,
	reflect.TypeOf(new(stringContaining)): 9,

	// multi-purpse matchers
	reflect.TypeOf(new(lengthOf)): 8,
	reflect.TypeOf(new(empty)):    7,

	// map & slice matchers
	reflect.TypeOf(new(keysContaining)):     6,
	reflect.TypeOf(new(elementsContaining)): 5,

	// type matchers
	reflect.TypeOf(new(implementerOf)):  4,
	reflect.TypeOf(new(convertibleTo)):  3,
	reflect.TypeOf(new(typeOf)):         2,
	reflect.TypeOf(new(anythingButNil)): 1,
	reflect.TypeOf(new(anything)):       0,
}
