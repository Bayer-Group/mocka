package mocka

import (
	"reflect"

	"github.com/Bayer-Group/mocka/v2/match"
)

// newCustomArguments constructor function for CustomArguments
func newCustomArguments(stub *Stub, arguments []interface{}) *CustomArguments {
	functionType := stub.toType()
	if isArgumentLengthValid(functionType, arguments) {
		reportInvalidArguments(stub.testReporter, functionType, arguments)
		return nil
	}

	matchers := getMatchers(functionType, arguments)
	if matchers == nil {
		reportInvalidArguments(stub.testReporter, functionType, arguments)
		return nil
	}

	return &CustomArguments{stub: stub, callCount: 0, argMatchers: matchers}
}

// isArgumentLengthValid returns whether or not the length of the provided arguments
// are valid for the stub. Taking in variadic function into account.
func isArgumentLengthValid(functionType reflect.Type, arguments []interface{}) bool {
	if functionType.IsVariadic() {
		return len(arguments) < functionType.NumIn()-1
	}

	return len(arguments) != functionType.NumIn()
}

// getMatchers returns a slice of matchers based on the types and values of the provided arguments
func getMatchers(functionType reflect.Type, arguments []interface{}) []match.SupportedKindsMatcher {
	matchers := make([]match.SupportedKindsMatcher, functionType.NumIn())
	for i := 0; i < functionType.NumIn(); i++ {
		aType := functionType.In(i)

		if isVariadicArgument(functionType, i) {
			if len(arguments) == functionType.NumIn()-1 {
				matchers[i] = match.Nil()
				return matchers
			}

			variadicArguments := arguments[i:]
			variadicMatchers := make([]match.SupportedKindsMatcher, len(variadicArguments))
			for sliceIndex, arg := range variadicArguments {
				m, found := getMatcher(arg, aType.Elem())
				if !found {
					return nil
				}

				variadicMatchers[sliceIndex] = m
			}

			matchers[i] = match.SliceOf(variadicMatchers...)
			return matchers
		}

		m, found := getMatcher(arguments[i], aType)
		if !found {
			return nil
		}

		matchers[i] = m
	}

	return matchers
}

// getMatcher returns a matcher for the provided type and value
func getMatcher(value interface{}, valueType reflect.Type) (match.SupportedKindsMatcher, bool) {
	if matcher, ok := value.(match.SupportedKindsMatcher); ok {
		if _, ok := matcher.SupportedKinds()[valueType.Kind()]; !ok {
			return nil, false
		}

		return matcher, true
	}

	if !areTypeAndValueEquivalent(valueType, value) {
		return nil, false
	}

	if value == nil {
		return match.Nil(), true
	}

	return match.Exactly(value), true
}

// CustomArguments represents a unique set of custom arguments in which
// the stubbed function will have different return values for
type CustomArguments struct {
	stub        *Stub
	argMatchers []match.SupportedKindsMatcher
	out         []interface{}
	onCalls     []*OnCall
	callCount   int
}

// Return sets the return values for this set of custom arguments
func (ca *CustomArguments) Return(returnValues ...interface{}) {
	ca.stub.lock.Lock()
	defer ca.stub.lock.Unlock()

	if !validateOutParameters(ca.stub.toType(), returnValues) {
		reportInvalidOutParameters(ca.stub.testReporter, ca.stub.toType(), returnValues)
		return
	}

	ca.out = returnValues
}

// OnCall returns an interface that allows for changing the
// return values based on the call index for this specific set
// of custom arguments.
func (ca *CustomArguments) OnCall(callIndex int) *OnCall {
	ca.stub.lock.Lock()
	defer ca.stub.lock.Unlock()

	for _, o := range ca.onCalls {
		if o.index == callIndex {
			return o
		}
	}

	o := &OnCall{index: callIndex, stub: ca.stub}
	ca.onCalls = append(ca.onCalls, o)
	return o
}

// OnFirstCall returns an interface that allows for changing the
// return values of the first call for this specific set
// of custom arguments.
func (ca *CustomArguments) OnFirstCall() *OnCall {
	return ca.OnCall(0)
}

// OnSecondCall returns an interface that allows for changing the
// return values of the second call for this specific set
// of custom arguments.
func (ca *CustomArguments) OnSecondCall() *OnCall {
	return ca.OnCall(1)
}

// OnThirdCall returns an interface that allows for changing the
// return values of the third call for this specific set
// of custom arguments.
func (ca *CustomArguments) OnThirdCall() *OnCall {
	return ca.OnCall(2)
}

// isMatch returns false if any of the argument matchers return false or
// if there is a panic from inside a matcher; otherwise true
func (ca *CustomArguments) isMatch(arguments []interface{}) (isMatch bool) {
	defer func() {
		if r := recover(); r != nil {
			isMatch = false
		}
	}()

	for i, arg := range arguments {
		if !ca.argMatchers[i].Match(arg) {
			return false
		}
	}

	return true
}
