package mocka

import (
	"fmt"
	"reflect"
	"strings"
)

// outParameterValidationError custom error for a out
// parameter validation error
type outParameterValidationError struct {
	fnType   reflect.Type
	provided []interface{}
}

// String returns a string that represents an out parameter validation error.
func (a *outParameterValidationError) String() string {
	if a.fnType == nil {
		return fmt.Sprintf("mocka: expected return values of (%v) to match function return values", strings.Join(mapToTypeName(a.provided), ", "))
	}

	real := make([]string, a.fnType.NumOut())
	for i := 0; i < a.fnType.NumOut(); i++ {
		real[i] = toFriendlyName(a.fnType.Out(i))
	}

	return fmt.Sprintf("mocka: expected return values of type (%v), but recieved (%v)", strings.Join(real, ", "), strings.Join(mapToTypeName(a.provided), ", "))
}

// Error returns a string that represents an out parameter validation error.
func (a *outParameterValidationError) Error() string {
	return a.String()
}
