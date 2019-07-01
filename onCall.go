package mocka

// OnCaller describes the functionality to set custom return value based on call index
type OnCaller interface {
	OnCall(int) Returner
	OnFirstCall() Returner
	OnSecondCall() Returner
	OnThirdCall() Returner
}

type onCall struct {
	stub  *mockFunction
	index int
	out   []interface{}
}

// Return sets the return values for this set of custom arguments
func (c *onCall) Return(returnValues ...interface{}) error {
	if !validateOutParameters(c.stub.toType(), returnValues) {
		return &outParameterValidationError{c.stub.toType(), returnValues}
	}

	c.out = returnValues
	return nil
}
