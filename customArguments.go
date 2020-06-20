package mocka

import (
	"errors"
	"reflect"

	"github.com/MonsantoCo/mocka/match"
)

// newCustomArguments constructor function for customArguments
func newCustomArguments(stub *mockFunction, arguments []interface{}) *customArguments {
	if stub == nil || stub.toType().Kind() != reflect.Func {
		return &customArguments{
			argValidationError: &argumentValidationError{
				provided: arguments,
			},
		}
	}

	functionType := stub.toType()
	if isArgumentLengthValid(functionType, arguments) {
		return &customArguments{
			argValidationError: &argumentValidationError{
				fnType:   functionType,
				provided: arguments,
			},
		}
	}

	matchers, err := getMatchers(functionType, arguments)
	return &customArguments{
		stub:               stub,
		callCount:          0,
		argMatchers:        matchers,
		argValidationError: err,
	}
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
func getMatchers(functionType reflect.Type, arguments []interface{}) ([]match.SupportedKindsMatcher, error) {
	matchers := make([]match.SupportedKindsMatcher, functionType.NumIn())
	for i := 0; i < functionType.NumIn(); i++ {
		aType := functionType.In(i)

		if isVariadicArgument(functionType, i) {
			if len(arguments) == functionType.NumIn()-1 {
				matchers[i] = match.Nil()
				return matchers, nil
			}

			variadicArguments := arguments[i:]
			variadicMatchers := make([]match.SupportedKindsMatcher, len(variadicArguments))
			for sliceIndex, arg := range variadicArguments {
				m, found := getMatcher(arg, aType.Elem())
				if !found {
					return nil, &argumentValidationError{
						fnType:   functionType,
						provided: arguments,
					}
				}

				variadicMatchers[sliceIndex] = m
			}

			matchers[i] = match.SliceOf(variadicMatchers...)
			return matchers, nil
		}

		m, found := getMatcher(arguments[i], aType)
		if !found {
			return nil, &argumentValidationError{
				fnType:   functionType,
				provided: arguments,
			}
		}

		matchers[i] = m
	}

	return matchers, nil
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

type customArguments struct {
	stub               *mockFunction
	argMatchers        []match.SupportedKindsMatcher
	argValidationError error
	out                []interface{}
	onCalls            []*onCall
	callCount          int
}

// Return sets the return values for this set of custom arguments
func (ca *customArguments) Return(returnValues ...interface{}) error {
	if ca.argValidationError != nil {
		return ca.argValidationError
	}

	if ca.stub == nil {
		return errors.New("mocka: stub does not exist")
	}

	ca.stub.lock.Lock()
	defer ca.stub.lock.Unlock()

	if !validateOutParameters(ca.stub.toType(), returnValues) {
		return &outParameterValidationError{ca.stub.toType(), returnValues}
	}

	ca.out = returnValues
	return nil
}

// OnCall returns an interface that allows for changing the
// return values based on the call index for this specific set
// of custom arguments.
func (ca *customArguments) OnCall(callIndex int) Returner {
	// TODO - future story
	// validate stub exists before using .lock
	// change return to also return an error if stub does not exist
	ca.stub.lock.Lock()
	defer ca.stub.lock.Unlock()

	for _, o := range ca.onCalls {
		if o.index == callIndex {
			return o
		}
	}

	o := &onCall{index: callIndex, stub: ca.stub}
	ca.onCalls = append(ca.onCalls, o)
	return o
}

// OnFirstCall returns an interface that allows for changing the
// return values of the first call for this specific set
// of custom arguments.
func (ca *customArguments) OnFirstCall() Returner {
	return ca.OnCall(0)
}

// OnSecondCall returns an interface that allows for changing the
// return values of the second call for this specific set
// of custom arguments.
func (ca *customArguments) OnSecondCall() Returner {
	return ca.OnCall(1)
}

// OnThirdCall returns an interface that allows for changing the
// return values of the third call for this specific set
// of custom arguments.
func (ca *customArguments) OnThirdCall() Returner {
	return ca.OnCall(2)
}

// isMatch returns false if any of the argument matchers return false or
// if there is a panic from inside a matcher; otherwise true
func (ca *customArguments) isMatch(arguments []interface{}) (isMatch bool) {
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
