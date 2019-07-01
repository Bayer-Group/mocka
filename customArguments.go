package mocka

type customArguments struct {
	stub      *mockFunction
	args      []interface{}
	out       []interface{}
	onCalls   []*onCall
	callCount int
}

// Return sets the return values for this set of custom arguments
func (ca *customArguments) Return(returnValues ...interface{}) error {
	if !validateOutParameters(ca.stub.toType(), returnValues) {
		return &outParameterValidationError{ca.stub.toType(), returnValues}
	}

	if !validateArguments(ca.stub.toType(), ca.args) {
		return &argumentValidationError{ca.stub.toType(), ca.args}
	}

	ca.out = returnValues
	return nil
}

// OnCall returns an interface that allows for changing the
// return values based on the call index for this specific set
// of custom arguments.
func (ca *customArguments) OnCall(callIndex int) Returner {
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
