package mocka

import (
	"fmt"
	"reflect"
	"strings"
)

// reportInvalidArguments reports invalid agument to fail the test
func reportInvalidArguments(testReporter TestReporter, functionType reflect.Type, arguments []interface{}) {
	real := make([]string, functionType.NumIn())
	for i := 0; i < functionType.NumIn(); i++ {
		if isVariadicArgument(functionType, i) {
			real[i] = fmt.Sprintf("...%v", toFriendlyName(functionType.In(i).Elem()))
			continue
		}

		real[i] = toFriendlyName(functionType.In(i))
	}

	testReporter.Errorf("mocka: expected arguments of type (%v), but received (%v)", strings.Join(real, ", "), strings.Join(mapToTypeName(arguments), ", "))
}

// reportInvalidOutParameters reports invalid out parameters to fail the test
func reportInvalidOutParameters(testReporter TestReporter, functionType reflect.Type, outParameters []interface{}) {
	if functionType == nil {
		testReporter.Errorf("mocka: expected return values of (%v) to match function return values", strings.Join(mapToTypeName(outParameters), ", "))
		return
	}

	real := make([]string, functionType.NumOut())
	for i := 0; i < functionType.NumOut(); i++ {
		real[i] = toFriendlyName(functionType.Out(i))
	}

	testReporter.Errorf("mocka: expected return values of type (%v), but received (%v)", strings.Join(real, ", "), strings.Join(mapToTypeName(outParameters), ", "))
}
