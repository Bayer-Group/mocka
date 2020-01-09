package mocka

import (
	"fmt"
	"reflect"
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
		return fmt.Sprintf("mocka: expected arguments of %v to match function arguments", mapToTypeName(a.provided))
	}

	real := make([]string, a.fnType.NumIn())
	for i := 0; i < a.fnType.NumIn(); i++ {
		real[i] = a.fnType.In(i).Name()
	}

	return fmt.Sprintf("mocka: expected arguments of type %v, but recieved %v", real, mapToTypeName(a.provided))
}

// Error returns a string that represents an argument validation error.
func (a *argumentValidationError) Error() string {
	return a.String()
}
