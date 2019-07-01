package mocka

import (
	"fmt"
	"reflect"
)

// outParameterValidationError custom error for a out
// parameter validation error
type outParameterValidationError struct {
	fnType   reflect.Type
	provided []interface{}
}

// Error returns a string that represents an out parameter validation error.
func (a *outParameterValidationError) Error() string {
	var real []string
	for i := 0; i < a.fnType.NumOut(); i++ {
		t := a.fnType.Out(i)
		real = append(real, t.Name())
	}

	return fmt.Sprintf("mocka: expected return values of type %v, but recieved %v", real, mapToTypeName(a.provided))
}
