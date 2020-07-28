// Package mocka is a simple mocking and stubbing library for the
// Go programming language. It is used to assist with writing
// unit tests around third-party functions.
//
// There are times when you would want to control the output of a
// third-party function in testing. Sometimes making an wrapper around
// that package/function is more than what you are wanting to do. Mocka
// is here to solve that problem. It allows you to control the output of
// functions without needing to write any additional code.
//
// Currently if you would want to control the output of a function in go
// it would be akin to
//
// 		// --- main.go ---
// 		// alias function for unit testing
// 		var jsonMarshal = json.Marshal
// 		...
// 		// --- main_test.go ---
// 		// create temporary variable to store original function
// 		var jsonMarshalOriginal func(v interface{}) ([]byte, error)
//		...
// 		func TestMarshal(t *testing.T) {
// 			jsonMarshalOriginal = jsonMarshal
// 			jsonMarshal = func(v interface{}) ([]byte, error) {
// 				return []byte("value"), nil
// 			}
// 			defer func() {
// 				jsonMarshal	= jsonMarshalOriginal
// 			}()
// 			// Your test code
// 		}
//
// This structure increases the length of unit tests; depending on how
// many functions are needing to control. Mocka provides a safe way to
// stub functions while also reducing the amount of code required.
//
// Mocka does this safely using reflection, no calls to the unsafe
// package are made.
//
// The mocka way would be
//
// 		// --- main.go ---
// 		// alias function for unit testing (in production code)
// 		var jsonMarshal = json.Marshal
// 		...
// 		// --- main_test.go ---
// 		func TestMarshal(t *testing.T) {
// 		    stub := mocka.Function(t, &jsonMarshal, []byte("value"), nil)
// 		    defer stub.Restore()
// 		    // Your test code
// 		}
//
package mocka

import "log"

// TestReporter is an interface used to fail tests.
// It is satisfied by the standard library testing.T and the
// response from GinkgoT()
type TestReporter interface {
	Errorf(string, ...interface{})
}

// Function replaces the provided function with a stubbed implementation. The
// stub has the ability to change change the return values of the original function
// in many different cases. The stub also provides the ability to get meta data
// associated to any call against the original function.
func Function(testReporter TestReporter, originalFuncPtr interface{}, returnValues ...interface{}) *Stub {
	return newStub(ensureTestReporter(testReporter, log.Fatal), originalFuncPtr, returnValues)
}

// CreateSandbox returns an isolated sandbox from which functions can be stubbed. The
// benefit you receive from using a sandbox is the ability to perform one call to Restore
// for a collection of Stubs
func CreateSandbox(testReporter TestReporter) *Sandbox {
	return &Sandbox{testReporter: ensureTestReporter(testReporter, log.Fatal)}
}

// ensureTestReporter returns the existing test reporter or calls exit
func ensureTestReporter(testReporter TestReporter, exit func(...interface{})) TestReporter {
	if testReporter == nil {
		exit("mocka: test reporter required to fail tests")
		return nil
	}
	return testReporter
}
