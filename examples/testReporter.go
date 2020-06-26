package examples

import (
	"log"
)

// t is used in place of *t.Testing for examples
var t = &mockTestReporter{}

// mockTestReporter used to simulate a test reporter
// to validate expected test failures
type mockTestReporter struct {
}

// Errorf appends the failure message to the internal messages slice
func (*mockTestReporter) Errorf(f string, args ...interface{}) {
	log.Fatalf(f, args...)
}
