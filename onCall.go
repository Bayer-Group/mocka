package mocka

// OnCall describes the functionality to set custom return value based on call index
type OnCall struct {
	stub  *Stub
	index int
	out   []interface{}
}

// Return sets the return values for this set of custom arguments
func (c *OnCall) Return(returnValues ...interface{}) {
	if !validateOutParameters(c.stub.toType(), returnValues) {
		c.stub.testReporter.Errorf("%v", &outParameterValidationError{c.stub.toType(), returnValues})
		return
	}

	c.stub.lock.Lock()
	defer c.stub.lock.Unlock()

	c.out = returnValues
}
