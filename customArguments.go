package mocka

import (
	"errors"
	"reflect"

	"github.com/MonsantoCo/mocka/match"
)

func newCustomArguments(stub *mockFunction, arguments []interface{}) *customArguments {
	if stub == nil || stub.toType().Kind() != reflect.Func {
		return &customArguments{
			argValidationError: &argumentValidationError{
				provided: arguments,
			},
		}
	}

	functionType := stub.toType()
	if len(arguments) != functionType.NumIn() {
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

func getMatchers(functionType reflect.Type, arguments []interface{}) ([]match.SupportedKindsMatcher, error) {
	matchers := make([]match.SupportedKindsMatcher, functionType.NumIn())
	for i, arg := range arguments {
		aType := functionType.In(i)

		switch arg.(type) {
		case match.SupportedKindsMatcher:
			matcher := arg.(match.SupportedKindsMatcher)
			if _, ok := matcher.SupportedKinds()[aType.Kind()]; !ok {
				return nil, &argumentValidationError{
					fnType:   functionType,
					provided: arguments,
				}
			}

			matchers[i] = matcher
		case nil:
			if !areTypeAndValueEquivalent(aType, arg) {
				return nil, &argumentValidationError{
					fnType:   functionType,
					provided: arguments,
				}
			}

			matchers[i] = match.Nil()
		default:
			if !areTypeAndValueEquivalent(aType, arg) {
				return nil, &argumentValidationError{
					fnType:   functionType,
					provided: arguments,
				}
			}

			matchers[i] = match.Exactly(arg)
		}
	}

	return matchers, nil
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
	if ca.stub == nil {
		return errors.New("mocka: stub does not exist")
	}

	if ca.argValidationError != nil {
		return ca.argValidationError
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

// isMatch returns false if any of the argument matcher return false or
// if there is a panic from inside a mather; otherwise true
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
