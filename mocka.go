// Package mocka provides a simple mocking and stubbing library to assist with
// writing unit tests.
//
// Currently if you would want to stub a function in go it would be akin to
//
//		// alias function for unit testing (in production code)
//		var jsonMarshal = json.Marshal
//		...
//		// create temporary variable to store original function (in unit test)
//		var jsonMarshalOriginal func Marshal(v interface{}) ([]byte, error)
//		...
//		jsonMarshalOriginal = jsonMarshal
//		jsonMarshal = func Marshal(v interface{}) ([]byte, error) {
//			return []byte("value"), nil
//		}
//		defer func() {
//			jsonMarshal	= jsonMarshalOriginal
//		}()
//
// This structure increases the length of unit tests; depending on how many
// functions are needing to be stubbed. Mocka provides a safe way to stub functions
// while also reducing the amount of code needed.
//
// For example
//
//		// alias function for unit testing (in production code)
//		var jsonMarshal = json.Marshal
//		...
//		// create stub to validate against
//		stub := mocka.Function(&jsonMarshal, []byte("value"), nil)
//		defer stub.Restore()
//
package mocka

import "log"

// TestReporter is an interface used fail tests.
// It is satisfied by the standard library testing.T and the
// response from GinkgoT()
type TestReporter interface {
	Errorf(string, ...interface{})
}

// Function replaces the provided function with a stubbed implementation. The
// stub has the ability to change change the return values of the original function
// in many different cases. The stub also provides the ability to get meta data
// associated to any call against the original function.
//
// Function also returns an error if the replacement of the original function
// with the stub failed.
func Function(testReporter TestReporter, originalFuncPtr interface{}, returnValues ...interface{}) *Stub {
	return newStub(ensureTestReporter(testReporter), originalFuncPtr, returnValues)
}

// CreateSandbox returns an isolated sandbox from which functions can be stubbed. The
// benefit you receive from using a sandbox is the ability to perform one call to Restore
// for a collection of Stubs
func CreateSandbox(testReporter TestReporter) *Sandbox {
	return &Sandbox{testReporter: ensureTestReporter(testReporter)}
}

// ensureTestReporter returns the existing test reporter or a new logger to Stderr
func ensureTestReporter(testReporter TestReporter) TestReporter {
	if testReporter == nil {
		log.Fatal("mocka: test reporter required to fail tests")
	}
	return testReporter
}
