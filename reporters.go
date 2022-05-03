package mocka

import (
	"fmt"
	"reflect"
	"strings"
)

// reportInvalidArguments reports invalid argument to fail the test
func reportInvalidArguments(testReporter TestReporter, functionType reflect.Type, arguments []interface{}) {
	realArgTypes := make([]string, functionType.NumIn())
	for i := 0; i < functionType.NumIn(); i++ {
		if isVariadicArgument(functionType, i) {
			realArgTypes[i] = fmt.Sprintf("...%v", toFriendlyName(functionType.In(i).Elem()))
			continue
		}

		realArgTypes[i] = toFriendlyName(functionType.In(i))
	}

	testReporter.Errorf("mocka: expected arguments of type (%v), but received (%v)", strings.Join(realArgTypes, ", "), strings.Join(mapToTypeName(arguments), ", "))
}

// reportInvalidOutParameters reports invalid out parameters to fail the test
func reportInvalidOutParameters(testReporter TestReporter, functionType reflect.Type, outParameters []interface{}) {
	if functionType == nil {
		testReporter.Errorf("mocka: expected return values of (%v) to match function return values", strings.Join(mapToTypeName(outParameters), ", "))
		return
	}

	realReturnTypes := make([]string, functionType.NumOut())
	for i := 0; i < functionType.NumOut(); i++ {
		realReturnTypes[i] = toFriendlyName(functionType.Out(i))
	}

	testReporter.Errorf("mocka: expected return values of type (%v), but received (%v)", strings.Join(realReturnTypes, ", "), strings.Join(mapToTypeName(outParameters), ", "))
}
