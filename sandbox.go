package mocka

import "sync"

// Sandbox describes an isolated environment that functions can be stubbed.
type Sandbox struct {
	lock sync.Mutex

	stubs []*Stub
}

// Function replaces the provided function with a stubbed implementation. The
// stub has the ability to change the return values of the original function
// in many different cases. The stub also provides the ability to get meta data
// associated to any call against the original function.
//
// Function also returns an error if the replacement of the original function
// with the stub failed.
func (s *Sandbox) Function(originalFuncPtr interface{}, returnValues ...interface{}) (*Stub, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	mockFn, err := newMockFunction(originalFuncPtr, returnValues)
	if err != nil {
		return nil, err
	}

	s.stubs = append(s.stubs, mockFn)

	return mockFn, nil
}

// Restore restores all the function stubs that were created via this sandbox to
// the original functionality they once held.
func (s *Sandbox) Restore() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, stub := range s.stubs {
		if stub != nil {
			stub.Restore()
		}
	}

	// clears out the slice to prevent a memory leak.
	s.stubs = nil
}
