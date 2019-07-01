package mocka

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("onCall", func() {
	var (
		fn     func(string, int) (int, error)
		mockFn *mockFunction
	)

	BeforeEach(func() {
		fn = func(str string, num int) (int, error) {
			return len(str) + num, nil
		}
		mockFn = &mockFunction{
			originalFunc:  nil,
			functionPtr:   &fn,
			outParameters: []interface{}{42, nil},
			execFunc:      func([]interface{}) {},
		}

	})

	Describe("Return", func() {
		It("returns an error if one out parameter type does not match", func() {
			ca := &onCall{
				stub:  mockFn,
				index: 0,
			}

			err := ca.Return(42, "nil")

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: expected return values of type [int error], but recieved [int string]"))
		})

		It("assigns the OutParameters and returns nil if everything is valid", func() {
			ca := &onCall{
				stub:  mockFn,
				index: 0,
			}

			err := ca.Return(42, nil)

			Expect(err).To(BeNil())
			Expect(ca.out).To(Equal([]interface{}{42, nil}))
		})
	})
})
