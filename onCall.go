package mocka

// OnCall describes the functionality to set custom return value based on call index
type OnCall struct {
	stub  *Stub
	index int
	out   []interface{}
}

// Return sets the return values for this set of custom arguments
func (c *OnCall) Return(returnValues ...interface{}) {
	c.stub.lock.Lock()
	defer c.stub.lock.Unlock()

	if !validateOutParameters(c.stub.toType(), returnValues) {
		reportInvalidOutParameters(c.stub.testReporter, c.stub.toType(), returnValues)
		return
	}

	c.out = returnValues
}
