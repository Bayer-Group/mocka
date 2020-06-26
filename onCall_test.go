package mocka

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OnCall", func() {
	var (
		fn   func(string, int) (int, error)
		stub *Stub
	)

	BeforeEach(func() {
		fn = func(str string, num int) (int, error) {
			return len(str) + num, nil
		}
		stub = &Stub{
			originalFunc:  nil,
			functionPtr:   &fn,
			outParameters: []interface{}{42, nil},
			execFunc:      func([]interface{}) {},
		}
	})

	Describe("Return", func() {
		It("returns an error if the stub is nil", func() {
			ca := &OnCall{
				index: 0,
			}

			err := ca.Return(42, "nil")

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("mocka: stub does not exist"))
		})

		It("returns an error if one out parameter type does not match", func() {
			ca := &OnCall{
				stub:  stub,
				index: 0,
			}

			err := ca.Return(42, "nil")

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("mocka: expected return values of type (int, error), but received (int, string)"))
		})

		It("assigns the OutParameters and returns nil if everything is valid", func() {
			ca := &OnCall{
				stub:  stub,
				index: 0,
			}

			err := ca.Return(42, nil)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(ca.out).To(Equal([]interface{}{42, nil}))
		})
	})
})
