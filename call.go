package mocka

// Call represents the information for a specific call invocation of the stubbed function
type Call struct {
	args []interface{}
	out  []interface{}
}

// Arguments returns the arguments that stub was called with.
func (c Call) Arguments() []interface{} {
	return c.args
}

// ReturnValues returns the return values that the stubbed implementation returned.
func (c Call) ReturnValues() []interface{} {
	return c.out
}
