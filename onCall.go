package mocka

import "errors"

// OnCall describes the functionality to set custom return value based on call index
type OnCall struct {
	stub  *Stub
	index int
	out   []interface{}
}

// Return sets the return values for this set of custom arguments
func (c *OnCall) Return(returnValues ...interface{}) error {
	if c.stub == nil {
		return errors.New("mocka: stub does not exist")
	}

	c.stub.lock.Lock()
	defer c.stub.lock.Unlock()

	if !validateOutParameters(c.stub.toType(), returnValues) {
		return &outParameterValidationError{c.stub.toType(), returnValues}
	}

	c.out = returnValues
	return nil
}
