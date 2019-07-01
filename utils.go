package mocka

import (
	"fmt"
	"reflect"
)

// mapToInterfaces maps a slice of reflection values to interface values
func mapToInterfaces(values []reflect.Value) []interface{} {
	interfaces := make([]interface{}, len(values))

	for index, value := range values {
		if value.IsValid() {
			interfaces[index] = value.Interface()
		} else {
			interfaces[index] = nil
		}
	}

	return interfaces
}

// mapToReflectValue maps a slice of interfaces to reflection values.
func mapToReflectValue(interfaces []interface{}) []reflect.Value {
	values := make([]reflect.Value, len(interfaces))
	for index, interfaceValue := range interfaces {
		values[index] = reflect.ValueOf(interfaceValue)
	}

	return values
}

// cloneValue creates a deep clone of a type and creates a new memory address
func cloneValue(source interface{}, destin interface{}) error {
	sourceValue := reflect.ValueOf(source)
	destinValue := reflect.ValueOf(destin)

	if sourceValue.Kind() != reflect.Ptr {
		return fmt.Errorf("mocka: expected source value for clone to be a pointer, but it was a %v", sourceValue.Kind().String())
	}

	if destinValue.Kind() != reflect.Ptr {
		return fmt.Errorf("mocka: expected destination value for clone to be a pointer, but it was a %v", destinValue.Kind().String())
	}

	sourceElem := sourceValue.Elem()
	newType := reflect.New(sourceElem.Type())
	newElem := newType.Elem()
	newElem.Set(sourceElem)
	destinValue.Elem().Set(newType.Elem())

	return nil
}

// validateArguments validates that the arguments provided match the argument types.
func validateArguments(functionType reflect.Type, arguments []interface{}) bool {
	if functionType.Kind() != reflect.Func {
		return false
	}

	argumentCount := functionType.NumIn()
	if len(arguments) != argumentCount {
		return false
	}

	isValid := true
	for i := 0; isValid && i < argumentCount; i++ {
		argumentType := functionType.In(i)
		isValid = areTypeAndValueEquivalent(argumentType, arguments[i])
	}

	return isValid
}

// validateOutParameters validates that the arguments provided match the argument types.
func validateOutParameters(functionType reflect.Type, outParameters []interface{}) bool {
	if functionType.Kind() != reflect.Func {
		return false
	}

	outParameterCount := functionType.NumOut()
	if len(outParameters) != outParameterCount {
		return false
	}

	isValid := true
	for i := 0; isValid && i < outParameterCount; i++ {
		outParameterType := functionType.Out(i)
		isValid = areTypeAndValueEquivalent(outParameterType, outParameters[i])
	}

	return isValid
}

// areTypeAndValueEquivalent does a kind wise check an interface{} to determine type equivalency
func areTypeAndValueEquivalent(originalType reflect.Type, val interface{}) bool {
	switch originalKind := originalType.Kind(); originalKind {
	case reflect.Func, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		if val == nil {
			return true
		}

		return originalKind == reflect.TypeOf(val).Kind()
	case reflect.Interface:
		if val == nil {
			return true
		}

		return reflect.TypeOf(val).Implements(originalType)
	default:
		return originalKind == reflect.TypeOf(val).Kind()
	}
}

// areEqual compares slices are interface values to see if they are equal
func areEqual(firsts, seconds []interface{}) (equal bool) {
	equal = true

	for i := 0; equal && i < len(firsts); i++ {
		equal = firsts[i] == seconds[i]
	}

	return
}

// mapToTypeName maps a slice of interface values to their type names
func mapToTypeName(interfaces []interface{}) (n []string) {
	for _, inter := range interfaces {
		if inter == nil {
			n = append(n, "<nil>")
		} else {
			t := reflect.TypeOf(inter)
			switch t.Kind() {
			case reflect.Ptr:
				n = append(n, "*"+t.Elem().Name())
			default:
				n = append(n, t.Name())
			}
		}
	}
	return
}
