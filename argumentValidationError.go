package mocka

import (
	"fmt"
	"reflect"
	"strings"
)

// argumentValidationError defines custom error for a argument
// validation error.
type argumentValidationError struct {
	fnType   reflect.Type
	provided []interface{}
}

// String returns a string that represents an argument validation error.
func (a *argumentValidationError) String() string {
	if a.fnType == nil {
		return fmt.Sprintf("mocka: expected arguments of (%v) to match function arguments", strings.Join(mapToTypeName(a.provided), ", "))
	}

	real := make([]string, a.fnType.NumIn())
	for i := 0; i < a.fnType.NumIn(); i++ {
		if isVariadicArgument(i, a.fnType) {
			real[i] = fmt.Sprintf("...%v", toFriendlyName(a.fnType.In(i).Elem()))
			continue
		}

		real[i] = toFriendlyName(a.fnType.In(i))
	}

	return fmt.Sprintf("mocka: expected arguments of type (%v), but received (%v)", strings.Join(real, ", "), strings.Join(mapToTypeName(a.provided), ", "))
}

// Error returns a string that represents an argument validation error.
func (a *argumentValidationError) Error() string {
	return a.String()
}
