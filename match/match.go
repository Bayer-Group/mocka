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
