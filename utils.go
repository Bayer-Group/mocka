package mocka

import (
	"fmt"
	"reflect"
	"strings"
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

// validateOutParameters validates that the arguments provided match the argument types.
func validateOutParameters(functionType reflect.Type, outParameters []interface{}) bool {
	if functionType == nil || functionType.Kind() != reflect.Func {
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
	if originalType == nil {
		return false
	}

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
		v := reflect.ValueOf(val)
		if !v.IsValid() {
			return false
		}

		return originalKind == v.Type().Kind()
	}
}

// mapToTypeName maps a slice of interface values to their type names
func mapToTypeName(interfaces []interface{}) []string {
	names := make([]string, len(interfaces))
	for i, value := range interfaces {
		names[i] = toFriendlyName(value)
	}

	return names
}

// toFriendlyName returns a type name is a more human readable string
func toFriendlyName(value interface{}) string {
	if value == nil {
		return "<nil>"
	}

	switch t := getType(value); t.Kind() {
	case reflect.Ptr:
		return "*" + toFriendlyName(t.Elem())
	case reflect.Slice:
		return fmt.Sprintf("[]%v", toFriendlyName(t.Elem()))
	case reflect.Array:
		return fmt.Sprintf("[%v]%v", t.Len(), toFriendlyName(t.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map[%v]%v", toFriendlyName(t.Key()), toFriendlyName(t.Elem()))
	case reflect.Chan:
		return toChannelFriendlyName(t)
	case reflect.Func:
		return toFunctionFriendlyName(t)
	default:
		return t.Name()
	}
}

// toChannelFriendlyName returns the friendly name for a channel
func toChannelFriendlyName(t reflect.Type) string {
	switch t.ChanDir() {
	case reflect.RecvDir:
		return fmt.Sprintf("<-chan %v", toFriendlyName(t.Elem()))
	case reflect.SendDir:
		return fmt.Sprintf("chan<- %v", toFriendlyName(t.Elem()))
	default:
		return fmt.Sprintf("chan %v", toFriendlyName(t.Elem()))
	}
}

// toFunctionFriendlyName returns the friendly name for a function
func toFunctionFriendlyName(t reflect.Type) string {
	args := make([]string, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		args[i] = toFriendlyName(t.In(i))
	}

	if t.NumOut() > 0 {
		out := make([]string, t.NumOut())
		for i := 0; i < t.NumOut(); i++ {
			out[i] = toFriendlyName(t.Out(i))
		}
		return fmt.Sprintf("func(%v) (%v) {}", strings.Join(args, ", "), strings.Join(out, ", "))
	}

	return fmt.Sprintf("func(%v) {}", strings.Join(args, ", "))
}

// getType returns the type of the argument
func getType(value interface{}) reflect.Type {
	switch value.(type) {
	case reflect.Type:
		return value.(reflect.Type)
	default:
		return reflect.TypeOf(value)
	}
}
