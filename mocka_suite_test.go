package mocka

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

func TestMocka(t *testing.T) {
	RegisterFailHandler(Fail)
	format.TruncatedDiff = false
	RunSpecs(t, "Mocka Testing Suite")
}

// mockTestReporter used to simulate a test reporter
// to validate expected test failures
type mockTestReporter struct {
	messages []string
}

// Errorf appends the failure message to the internal messages slice
func (m *mockTestReporter) Errorf(f string, args ...interface{}) {
	m.messages = append(m.messages, fmt.Sprintf(f, args...))
}
