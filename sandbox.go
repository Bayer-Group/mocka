package mocka

// Sandbox describes an isolated environment that functions can be stubbed.
type Sandbox interface {
	StubFunction(interface{}, ...interface{}) (Stub, error)
	Restore()
}

type sandbox struct {
	stubs []mockFunction
}

// StubFunction replaces the provided function with a stubbed implementation. The
// stub has the ability to change the return values of the original function
// in many different cases. The stub also provides the ability to get meta data
// associated to any call against the original function.
//
// StubFunction also returns an error if the replacement of the original function
// with the stub failed.
func (s *sandbox) StubFunction(originalFuncPtr interface{}, returnValues ...interface{}) (Stub, error) {
	mockFn, err := newMockFunction(originalFuncPtr, returnValues)
	if err != nil {
		return nil, err
	}

	s.stubs = append(s.stubs, *mockFn)

	return mockFn, nil
}

// Restore restores all the function stubs that were created via this sandbox to
// the original functionality they once held.
func (s *sandbox) Restore() {
	for _, stub := range s.stubs {
		stub.Restore()
	}

	// clears out the slice to prevent a memory leak.
	s.stubs = nil
}
