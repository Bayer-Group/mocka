package mocka

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMocka(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mocka Testing Suite")
}
