package mocka

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OnCall", func() {
	var (
		fn               func(string, int) (int, error)
		stub             *Stub
		failTestReporter *mockTestReporter
	)

	BeforeEach(func() {
		fn = func(str string, num int) (int, error) {
			return len(str) + num, nil
		}
		stub = &Stub{
			testReporter:  GinkgoT(),
			originalFunc:  nil,
			functionPtr:   &fn,
			outParameters: []interface{}{42, nil},
			execFunc:      func([]interface{}) {},
		}
		failTestReporter = &mockTestReporter{}
	})

	Describe("Return", func() {
		It("returns an error if one out parameter type does not match", func() {
			stub.testReporter = failTestReporter
			ca := &OnCall{
				stub:  stub,
				index: 0,
			}

			ca.Return(42, "nil")

			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected return values of type (int, error), but received (int, string)",
			}))
		})

		It("assigns the OutParameters and returns nil if everything is valid", func() {
			ca := &OnCall{
				stub:  stub,
				index: 0,
			}

			ca.Return(42, nil)

			Expect(ca.out).To(Equal([]interface{}{42, nil}))
		})
	})
})
