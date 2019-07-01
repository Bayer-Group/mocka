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

// Error returns a string that represents an argument validation error.
func (a *argumentValidationError) Error() string {
	var real []string
	for i := 0; i < a.fnType.NumIn(); i++ {
		t := a.fnType.In(i)
		real = append(real, t.Name())
	}

	return fmt.Sprintf("mocka: expected arguments of type %v, but recieved %v", real, mapToTypeName(a.provided))
}
