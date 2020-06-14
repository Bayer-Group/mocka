package mocka

import (
	"reflect"
)

// Call describes a specific call the to the stub function
type Call interface {
	Arguments() []interface{}
	ReturnValues() []interface{}
}

// newCall returns a new call where the arguments are spread for variadic functions
func newCall(functionType reflect.Type, args []interface{}, out []interface{}) call {
	if functionType.IsVariadic() {
		var vArgs []interface{}

		for argIndex, arg := range args {
			if isVariadicArgument(argIndex, functionType) {
				slice := reflect.ValueOf(arg)
				for sliceIndex := 0; sliceIndex < slice.Len(); sliceIndex++ {
					vArgs = append(vArgs, slice.Index(sliceIndex).Interface())
				}
				continue
			}

			vArgs = append(vArgs, arg)
		}

		return call{args: vArgs, out: out}
	}

	return call{args: args, out: out}
}

type call struct {
	args []interface{}
	out  []interface{}
}

// Arguments returns the arguments that stub was called with.
func (c call) Arguments() []interface{} {
	// TODO - callout if a function is variadic should the
	// argument here be the variadic ones or the one that
	// has the variadic list?
	return c.args
}

// ReturnValues returns the return values that the stubbed implementation returned.
func (c call) ReturnValues() []interface{} {
	return c.out
}
