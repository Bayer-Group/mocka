// Package match provides matchers to use with mocka for stubbing functions with
// custom arguments.
package match

import "reflect"

// SupportedKindsMatcher describes the functionality of a custom argument matcher for mocka
type SupportedKindsMatcher interface {
	// SupportedKinds returns the supported kinds for the matcher
	SupportedKinds() map[reflect.Kind]struct{}

	// Match return true is the match was successful; otherwise false
	Match(interface{}) bool
}

// Priority returns the matchers priority to be compared against
func Priority(m SupportedKindsMatcher) int {
	switch m.(type) {
	case *anything:
		return 0
	case *anythingButNil:
		return 0
	case *convertibleTo:
		return 0
	case *elementsContaining:
		return 0
	case *empty:
		return 0
	case *exactly:
		return 10
	case *floatGreaterThan:
		return 0
	case *floatGreaterThanOrEqualTo:
		return 0
	case *floatLessThan:
		return 0
	case *floatLessThanOrEqualTo:
		return 0
	case *implementerOf:
		return 0
	case *intGreaterThan:
		return 0
	case *intGreaterThanOrEqualTo:
		return 0
	case *intLessThan:
		return 0
	case *intLessThanOrEqualTo:
		return 0
	case *keysContaining:
		return 0
	case *lengthOf:
		return 0
	case *nilMatcher:
		return 0
	case *stringContaining:
		return 0
	case *stringPrefix:
		return 0
	case *stringSuffix:
		return 0
	case *typeOf:
		return 0
	default:
		return 0
	}
}
