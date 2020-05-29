package mocka

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/MonsantoCo/mocka/match"
	"github.com/pkg/errors"
)

// variables used for unit testing
var _cloneValue = cloneValue

// Returner describes the functionality to update the return values to the stub
type Returner interface {
	Return(...interface{}) error
}

// GetCaller describes the functionality to get information for a specific call to the stub
type GetCaller interface {
	GetCalls() []Call
	GetCall(int) Call
	GetFirstCall() Call
	GetSecondCall() Call
	GetThirdCall() Call
	CallCount() int
	CalledOnce() bool
	CalledTwice() bool
	CalledThrice() bool
}

// OnCallReturner describes the functionality to update the return values itself of based on the call index
type OnCallReturner interface {
	OnCaller
	Returner
}

// Stub describes the functionality related to the stub replacement of a function
type Stub interface {
	Returner
	GetCaller
	OnCaller
	WithArgs(...interface{}) OnCallReturner
	ExecOnCall(func([]interface{}))
	Restore()
}

type mockFunction struct {
	lock sync.RWMutex

	originalFunc  interface{}
	functionPtr   interface{}
	outParameters []interface{}
	calls         []call
	customArgs    []*customArguments
	onCalls       []*onCall
	execFunc      func([]interface{})
}

// newMockFunction creates a new mock function and overrides the implementation of the original function.
func newMockFunction(originalFuncPtr interface{}, returnValues []interface{}) (*mockFunction, error) {
	if originalFuncPtr == nil {
		return nil, errors.New("mocka: expected the first argument to be a pointer to a function, but received a nil")
	}

	originalFuncValue := reflect.ValueOf(originalFuncPtr)
	if originalFuncValue.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("mocka: expected the first argument to be a pointer to a function, but received a %v", originalFuncValue.Kind().String())
	}

	originalFunc := originalFuncValue.Elem()
	if originalFunc.Kind() != reflect.Func {
		return nil, fmt.Errorf("mocka: expected the first argument to be a pointer to a function, but received a pointer to a %v", originalFunc.Kind().String())
	}

	if !validateOutParameters(originalFunc.Type(), returnValues) {
		return nil, &outParameterValidationError{originalFunc.Type(), returnValues}
	}

	stub := &mockFunction{
		originalFunc:  nil,
		functionPtr:   originalFuncPtr,
		outParameters: returnValues,
		execFunc:      func([]interface{}) {},
	}

	// Need to perform a deep clone to get a new pointer and memory address
	err := _cloneValue(originalFuncPtr, &stub.originalFunc)
	if err != nil {
		return nil, errors.Wrap(err, "mocka: could not clone function pointer to new memory address")
	}

	// Replace the original function the mock function implementation
	originalType := originalFunc.Type()
	originalFunc.Set(reflect.MakeFunc(originalType, stub.implementation))

	return stub, nil
}

// toType gets the reflection type from the mock function pointer
func (mf *mockFunction) toType() reflect.Type {
	return reflect.ValueOf(mf.functionPtr).Elem().Type()
}

// implementation defines the function that replaces the original
// function's functionality
func (mf *mockFunction) implementation(arguments []reflect.Value) []reflect.Value {
	mf.lock.Lock()
	defer mf.lock.Unlock()

	functionType := mf.toType()
	argumentsAsInterfaces := mapToInterfaces(arguments)
	outParameters, maybeCustomArguments := mf.getReturnValues(argumentsAsInterfaces, functionType)
	outParametersAsValues := mapToReflectValue(outParameters)

	outParametersAsInterfaces := make([]interface{}, len(outParametersAsValues))
	for index, value := range outParametersAsValues {
		outParamType := functionType.Out(index)
		if value.IsValid() {
			if outParamType.Kind() == reflect.Interface {
				newType := reflect.New(outParamType)
				newElem := newType.Elem()
				newElem.Set(reflect.ValueOf(value.Interface()))
				outParametersAsValues[index] = newElem
			}
			outParametersAsInterfaces[index] = value.Interface()
		} else {
			outParametersAsValues[index] = reflect.Zero(outParamType)
			outParametersAsInterfaces[index] = nil
		}
	}

	mf.execFunc(argumentsAsInterfaces)

	mf.calls = append(mf.calls, call{
		args: argumentsAsInterfaces,
		out:  outParametersAsInterfaces,
	})

	if maybeCustomArguments != nil {
		maybeCustomArguments.callCount++
	}

	return outParametersAsValues
}

// getReturnValues returns the correct out parameters based on the
// arguments passed into the function.
//
// This function also takes into account the current call index of function.
func (mf *mockFunction) getReturnValues(arguments []interface{}, functionType reflect.Type) ([]interface{}, *customArguments) {
	out := mf.outParameters

	for _, o := range mf.onCalls {
		if o.index == len(mf.calls) && o.out != nil {
			out = o.out
			break
		}
	}

	maybeCustomArgs := getHighestPriority(getPossible(mf.customArgs, arguments), functionType.NumIn())
	if maybeCustomArgs == nil {
		return out, nil
	}

	out = maybeCustomArgs.out
	for _, o := range maybeCustomArgs.onCalls {
		if o.index == maybeCustomArgs.callCount && o.out != nil {
			return o.out, maybeCustomArgs
		}
	}

	return out, maybeCustomArgs
}

// getHighestPriority returns the highest priority custom arguments if found;
// otherwise a nil
func getHighestPriority(customArgs []*customArguments, numArgs int) *customArguments {
	switch len(customArgs) {
	case 0:
		return nil
	case 1:
		return customArgs[0]
	}

	for i := 0; i < numArgs; i++ {
		highestPriority := 0
		newCustomArgs := make([]*customArguments, 0)

		for _, ca := range customArgs {
			if p := match.Priority(ca.argMatchers[i]); p > highestPriority {
				highestPriority = p
				newCustomArgs = make([]*customArguments, 0)
			}
			newCustomArgs = append(newCustomArgs, ca)
		}

		switch len(newCustomArgs) {
		case 0:
			return nil
		case 1:
			return newCustomArgs[0]
		default:
			customArgs = newCustomArgs
		}
	}

	return customArgs[0]
}

// getPossible returns the possible custom arguments
// that match the provided arguments
func getPossible(customArgs []*customArguments, arguments []interface{}) (possible []*customArguments) {
	for _, ca := range customArgs {
		if ca != nil && ca.out != nil && ca.match(arguments) {
			possible = append(possible, ca)
		}
	}
	return
}

// Returns updates the default out parameters returned when
// the mock function is called
func (mf *mockFunction) Return(returnValues ...interface{}) error {
	mf.lock.Lock()
	defer mf.lock.Unlock()

	if !validateOutParameters(mf.toType(), returnValues) {
		return &outParameterValidationError{mf.toType(), returnValues}
	}

	mf.outParameters = returnValues
	return nil
}

// WithArgs returns a StubWithArgs that can change the out parameters
// returned based on the arguments provided to this function
func (mf *mockFunction) WithArgs(arguments ...interface{}) OnCallReturner {
	mf.lock.Lock()
	defer mf.lock.Unlock()

	newCA := newCustomArguments(mf, arguments)
	for _, ca := range mf.customArgs {
		if ca != nil && ca.matchersAreEqual(newCA) {
			return ca
		}
	}

	mf.customArgs = append(mf.customArgs, newCA)

	return newCA
}

// CallCount returns the number of times the original function was called
// after the function was stubbed
func (mf *mockFunction) CallCount() int {
	mf.lock.RLock()
	defer mf.lock.RUnlock()

	return len(mf.calls)
}

// GetCalls returns all calls made to the original function that were
// captured by the stubbed implementation
func (mf *mockFunction) GetCalls() []Call {
	mf.lock.RLock()
	defer mf.lock.RUnlock()

	calls := make([]Call, len(mf.calls))

	for i, call := range mf.calls {
		c := call
		calls[i] = &c
	}

	return calls
}

// GetCall returns the arguments and return values of the original function
// that was captured by the stubbed implementation. It will return these values
// for the specified time the function was called.
//
// GetCall will also panic if the call index is lower than zero or greater than
// the number of times the function was called.
//
// The call index uses zero-based indexing
func (mf *mockFunction) GetCall(callIndex int) Call {
	mf.lock.RLock()
	defer mf.lock.RUnlock()

	if callIndex < 0 || callIndex >= mf.CallCount() {
		panic(fmt.Errorf("mocka: attempted to get CallMetaData for call %v, when the function has only been called %v times", callIndex, len(mf.calls)))
	}

	return &mf.calls[callIndex]
}

// GetFirstCall returns the arguments and return values of the original function
// that was captured by the stubbed implementation for the first time the original
// function was called.
//
// GetFirstCall will also panic if the original function was not called at least
// once.
func (mf *mockFunction) GetFirstCall() Call {
	return mf.GetCall(0)
}

// GetSecondCall returns the arguments and return values of the original function
// that was captured by the stubbed implementation for the second time the original
// function was called.
//
// GetSecondCall will also panic if the original function was not called at least
// twice.
func (mf *mockFunction) GetSecondCall() Call {
	return mf.GetCall(1)
}

// GetThirdCall returns the arguments and return values of the original function
// that was captured by the stubbed implementation for the three time the original
// function was called.
//
// GetThirdCall will also panic if the original function was not called at least
// thrice.
func (mf *mockFunction) GetThirdCall() Call {
	return mf.GetCall(2)
}

// CalledOnce returns true if the original function has been called at least once;
// otherwise, it will return false.
func (mf *mockFunction) CalledOnce() bool {
	return mf.CallCount() >= 1
}

// CalledTwice returns true if the original function has been called at twice once;
// otherwise, it will return false.
func (mf *mockFunction) CalledTwice() bool {
	return mf.CallCount() >= 2
}

// CalledThrice returns true if the original function has been called at thrice once;
// otherwise, it will return false.
func (mf *mockFunction) CalledThrice() bool {
	return mf.CallCount() >= 3
}

// OnCall returns an interface that allows for changing the
// return values based on the call index.
func (mf *mockFunction) OnCall(index int) Returner {
	mf.lock.Lock()
	defer mf.lock.Unlock()

	for _, o := range mf.onCalls {
		if o.index == index {
			return o
		}
	}

	o := &onCall{index: index, stub: mf}
	mf.onCalls = append(mf.onCalls, o)
	return o
}

// OnFirstCall returns an interface that allows for changing the
// return values of the first call.
func (mf *mockFunction) OnFirstCall() Returner {
	return mf.OnCall(0)
}

// OnSecondCall returns an interface that allows for changing the
// return values of the second call.
func (mf *mockFunction) OnSecondCall() Returner {
	return mf.OnCall(1)
}

// OnThirdCall returns an interface that allows for changing the
// return values of the third call.
func (mf *mockFunction) OnThirdCall() Returner {
	return mf.OnCall(2)
}

// Restore removes the stub and restores the the original
// functionality back to the method
func (mf *mockFunction) Restore() {
	mf.lock.Lock()
	defer mf.lock.Unlock()

	valueOforiginalFunc := reflect.ValueOf(mf.originalFunc)
	functionValue := reflect.ValueOf(mf.functionPtr).Elem()

	functionValue.Set(valueOforiginalFunc)
}

// ExecOnCall assigns a function to be called when the stub
// implementation is called.
func (mf *mockFunction) ExecOnCall(execFunc func([]interface{})) {
	mf.lock.Lock()
	defer mf.lock.Unlock()

	mf.execFunc = execFunc
}
