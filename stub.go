package mocka

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/MonsantoCo/mocka/match"
)

// variables used for unit testing
var _cloneValue = cloneValue

// Stub represents the stub for a function
type Stub struct {
	lock sync.RWMutex

	testReporter  TestReporter
	originalFunc  interface{}
	functionPtr   interface{}
	outParameters []interface{}
	calls         []Call
	customArgs    []*CustomArguments
	onCalls       []*OnCall
	execFunc      func([]interface{})
}

// newStub creates a stub function and overrides the implementation of the original function.
func newStub(testReporter TestReporter, originalFuncPtr interface{}, returnValues []interface{}) *Stub {
	if originalFuncPtr == nil {
		testReporter.Errorf("mocka: expected the first argument to be a pointer to a function, but received a nil")
		return nil
	}

	originalFuncValue := reflect.ValueOf(originalFuncPtr)
	if originalFuncValue.Kind() != reflect.Ptr {
		testReporter.Errorf("mocka: expected the first argument to be a pointer to a function, but received a %v", originalFuncValue.Kind().String())
		return nil
	}

	originalFunc := originalFuncValue.Elem()
	if originalFunc.Kind() != reflect.Func {
		testReporter.Errorf("mocka: expected the first argument to be a pointer to a function, but received a pointer to a %v", originalFunc.Kind().String())
		return nil
	}

	if !validateOutParameters(originalFunc.Type(), returnValues) {
		testReporter.Errorf("%v", &outParameterValidationError{originalFunc.Type(), returnValues})
		return nil
	}

	stub := &Stub{
		originalFunc:  nil,
		testReporter:  testReporter,
		functionPtr:   originalFuncPtr,
		outParameters: returnValues,
		execFunc:      func([]interface{}) {},
	}

	// Need to perform a deep clone to get a new pointer and memory address
	err := _cloneValue(originalFuncPtr, &stub.originalFunc)
	if err != nil {
		stub.testReporter.Errorf("mocka: could not clone function pointer to new memory address: %v", err)
		return nil
	}

	// Replace the original function the mock function implementation
	originalType := originalFunc.Type()
	originalFunc.Set(reflect.MakeFunc(originalType, stub.implementation))

	return stub
}

// toType gets the reflection type from the mock function pointer
func (stub *Stub) toType() reflect.Type {
	return reflect.ValueOf(stub.functionPtr).Elem().Type()
}

// implementation defines the function that replaces the original
// function's functionality
func (stub *Stub) implementation(arguments []reflect.Value) []reflect.Value {
	stub.lock.Lock()
	defer stub.lock.Unlock()

	functionType := stub.toType()
	argumentsAsInterfaces := mapToInterfaces(arguments)
	outParameters, maybeCustomArguments := stub.getReturnValues(argumentsAsInterfaces, functionType)
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

	stub.execFunc(argumentsAsInterfaces)

	stub.calls = append(stub.calls, Call{args: argumentsAsInterfaces, out: outParametersAsInterfaces})

	if maybeCustomArguments != nil {
		maybeCustomArguments.callCount++
	}

	return outParametersAsValues
}

// getReturnValues returns the correct out parameters based on the
// arguments passed into the function.
//
// This function also takes into account the current call index of function.
func (stub *Stub) getReturnValues(arguments []interface{}, functionType reflect.Type) ([]interface{}, *CustomArguments) {
	out := stub.outParameters

	for _, o := range stub.onCalls {
		if o.index == len(stub.calls) && o.out != nil {
			out = o.out
			break
		}
	}

	maybeCustomArgs := getHighestPriority(getPossible(stub.customArgs, arguments), functionType.NumIn())
	if maybeCustomArgs == nil {
		return out, nil
	}

	if maybeCustomArgs.out != nil {
		out = maybeCustomArgs.out
	}

	for _, o := range maybeCustomArgs.onCalls {
		if o.index == maybeCustomArgs.callCount && o.out != nil {
			return o.out, maybeCustomArgs
		}
	}

	return out, maybeCustomArgs
}

// getHighestPriority returns the highest priority custom arguments if found;
// otherwise a nil
func getHighestPriority(customArgs []*CustomArguments, numArgs int) *CustomArguments {
	switch len(customArgs) {
	case 0:
		return nil
	case 1:
		return customArgs[0]
	}

	for i := 0; i < numArgs; i++ {
		var highestPriority float64
		newCustomArgs := make([]*CustomArguments, 0)

		for _, ca := range customArgs {
			p := match.Priority(ca.argMatchers[i])
			if p > highestPriority {
				highestPriority = p
				newCustomArgs = make([]*CustomArguments, 0)
			}

			if p == highestPriority {
				newCustomArgs = append(newCustomArgs, ca)
			}
		}

		if len(newCustomArgs) == 1 {
			return newCustomArgs[0]
		}

		customArgs = newCustomArgs
	}

	return customArgs[0]
}

// getPossible returns the possible custom arguments
// that match the provided arguments
func getPossible(customArgs []*CustomArguments, arguments []interface{}) (possible []*CustomArguments) {
	for _, ca := range customArgs {
		if ca != nil && ca.isMatch(arguments) {
			possible = append(possible, ca)
		}
	}
	return
}

// Return updates the default out parameters returned when
// the mock function is called
func (stub *Stub) Return(returnValues ...interface{}) error {
	stub.lock.Lock()
	defer stub.lock.Unlock()

	if !validateOutParameters(stub.toType(), returnValues) {
		return &outParameterValidationError{stub.toType(), returnValues}
	}

	stub.outParameters = returnValues
	return nil
}

// WithArgs returns a StubWithArgs that can change the out parameters
// returned based on the arguments provided to this function
func (stub *Stub) WithArgs(arguments ...interface{}) *CustomArguments {
	stub.lock.Lock()
	defer stub.lock.Unlock()

	newCA := newCustomArguments(stub, arguments)
	for _, ca := range stub.customArgs {
		if ca != nil && reflect.DeepEqual(ca.argMatchers, newCA.argMatchers) {
			return ca
		}
	}

	stub.customArgs = append(stub.customArgs, newCA)

	return newCA
}

// CallCount returns the number of times the original function was called
// after the function was stubbed
func (stub *Stub) CallCount() int {
	stub.lock.RLock()
	defer stub.lock.RUnlock()

	return len(stub.calls)
}

// GetCalls returns all calls made to the original function that were
// captured by the stubbed implementation
func (stub *Stub) GetCalls() []Call {
	stub.lock.RLock()
	defer stub.lock.RUnlock()

	calls := make([]Call, len(stub.calls))

	for i, call := range stub.calls {
		c := call
		calls[i] = c
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
func (stub *Stub) GetCall(callIndex int) Call {
	stub.lock.RLock()
	defer stub.lock.RUnlock()

	if callIndex < 0 || callIndex >= stub.CallCount() {
		panic(fmt.Errorf("mocka: attempted to get CallMetaData for call %v, when the function has only been called %v times", callIndex, len(stub.calls)))
	}

	return stub.calls[callIndex]
}

// GetFirstCall returns the arguments and return values of the original function
// that was captured by the stubbed implementation for the first time the original
// function was called.
//
// GetFirstCall will also panic if the original function was not called at least
// once.
func (stub *Stub) GetFirstCall() Call {
	return stub.GetCall(0)
}

// GetSecondCall returns the arguments and return values of the original function
// that was captured by the stubbed implementation for the second time the original
// function was called.
//
// GetSecondCall will also panic if the original function was not called at least
// twice.
func (stub *Stub) GetSecondCall() Call {
	return stub.GetCall(1)
}

// GetThirdCall returns the arguments and return values of the original function
// that was captured by the stubbed implementation for the three time the original
// function was called.
//
// GetThirdCall will also panic if the original function was not called at least
// thrice.
func (stub *Stub) GetThirdCall() Call {
	return stub.GetCall(2)
}

// CalledOnce returns true if the original function has been called at least once;
// otherwise, it will return false.
func (stub *Stub) CalledOnce() bool {
	return stub.CallCount() >= 1
}

// CalledTwice returns true if the original function has been called at twice once;
// otherwise, it will return false.
func (stub *Stub) CalledTwice() bool {
	return stub.CallCount() >= 2
}

// CalledThrice returns true if the original function has been called at thrice once;
// otherwise, it will return false.
func (stub *Stub) CalledThrice() bool {
	return stub.CallCount() >= 3
}

// OnCall returns an interface that allows for changing the
// return values based on the call index.
func (stub *Stub) OnCall(index int) *OnCall {
	stub.lock.Lock()
	defer stub.lock.Unlock()

	for _, o := range stub.onCalls {
		if o.index == index {
			return o
		}
	}

	o := &OnCall{index: index, stub: stub}
	stub.onCalls = append(stub.onCalls, o)
	return o
}

// OnFirstCall returns an interface that allows for changing the
// return values of the first call.
func (stub *Stub) OnFirstCall() *OnCall {
	return stub.OnCall(0)
}

// OnSecondCall returns an interface that allows for changing the
// return values of the second call.
func (stub *Stub) OnSecondCall() *OnCall {
	return stub.OnCall(1)
}

// OnThirdCall returns an interface that allows for changing the
// return values of the third call.
func (stub *Stub) OnThirdCall() *OnCall {
	return stub.OnCall(2)
}

// Restore removes the stub and restores the the original
// functionality back to the method
func (stub *Stub) Restore() {
	stub.lock.Lock()
	defer stub.lock.Unlock()

	valueOforiginalFunc := reflect.ValueOf(stub.originalFunc)
	functionValue := reflect.ValueOf(stub.functionPtr).Elem()

	functionValue.Set(valueOforiginalFunc)
}

// ExecOnCall assigns a function to be called when the stub
// implementation is called.
func (stub *Stub) ExecOnCall(execFunc func([]interface{})) {
	stub.lock.Lock()
	defer stub.lock.Unlock()

	stub.execFunc = execFunc
}
