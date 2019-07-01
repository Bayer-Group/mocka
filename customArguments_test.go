package mocka

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("customArguments", func() {
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
			ca := &customArguments{
				stub: mockFn,
				args: []interface{}{"", 42},
			}

			err := ca.Return(42, "nil")

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: expected return values of type [int error], but recieved [int string]"))
		})

		It("returns an error if one argument type does not match", func() {
			ca := &customArguments{
				stub: mockFn,
				args: []interface{}{"", "42"},
			}

			err := ca.Return(42, nil)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: expected arguments of type [string int], but recieved [string string]"))
		})

		It("assigns the OutParameters and returns nil if everything is valid", func() {
			ca := &customArguments{
				stub: mockFn,
				args: []interface{}{"", 42},
			}

			err := ca.Return(42, nil)

			Expect(err).To(BeNil())
			Expect(ca.out).To(Equal([]interface{}{42, nil}))
		})
	})

	Describe("OnCall", func() {
		var ca *customArguments

		BeforeEach(func() {
			ca = &customArguments{
				stub: mockFn,
				onCalls: []*onCall{
					&onCall{stub: mockFn, index: 1},
					&onCall{stub: mockFn, index: 2},
					&onCall{stub: mockFn, index: 0},
				},
			}
		})

		It("returns a pointer to an onCall struct", func() {
			result := ca.OnCall(5)

			o, ok := result.(*onCall)

			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockFn, index: 5}))
		})

		It("appends the new onCall struct to the onCalls slice", func() {
			_ = ca.OnCall(5)

			Expect(ca.onCalls).To(HaveLen(4))

		})

		It("returns an existing onCall object if one exists for that index", func() {
			result := ca.OnCall(2)

			o, ok := result.(*onCall)

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockFn, index: 2}))

		})
	})

	Describe("OnFirstCall", func() {
		It("creates a new onCall with a 0 index", func() {
			ca := &customArguments{
				stub: mockFn,
				onCalls: []*onCall{
					&onCall{stub: mockFn, index: 1},
					&onCall{stub: mockFn, index: 2},
				},
			}

			result := ca.OnFirstCall()

			o, ok := result.(*onCall)

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockFn, index: 0}))

		})

		It("returns the existing first call", func() {
			ca := &customArguments{
				stub: mockFn,
				onCalls: []*onCall{
					&onCall{stub: mockFn, index: 1},
					&onCall{stub: mockFn, index: 2},
					&onCall{stub: mockFn, index: 0},
				},
			}

			result := ca.OnFirstCall()

			o, ok := result.(*onCall)

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockFn, index: 0}))

		})
	})

	Describe("OnSecondCall", func() {
		It("creates a new onCall with a 1 index", func() {
			ca := &customArguments{
				stub: mockFn,
				onCalls: []*onCall{
					&onCall{stub: mockFn, index: 0},
					&onCall{stub: mockFn, index: 2},
				},
			}

			result := ca.OnSecondCall()

			o, ok := result.(*onCall)

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockFn, index: 1}))

		})

		It("returns the existing second call", func() {
			ca := &customArguments{
				stub: mockFn,
				onCalls: []*onCall{
					&onCall{stub: mockFn, index: 1},
					&onCall{stub: mockFn, index: 2},
					&onCall{stub: mockFn, index: 0},
				},
			}

			result := ca.OnSecondCall()

			o, ok := result.(*onCall)

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockFn, index: 1}))

		})
	})

	Describe("OnThirdCall", func() {
		It("creates a new onCall with a 2 index", func() {
			ca := &customArguments{
				stub: mockFn,
				onCalls: []*onCall{
					&onCall{stub: mockFn, index: 0},
					&onCall{stub: mockFn, index: 1},
				},
			}

			result := ca.OnThirdCall()

			o, ok := result.(*onCall)

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockFn, index: 2}))

		})

		It("returns the existing third call", func() {
			ca := &customArguments{
				stub: mockFn,
				onCalls: []*onCall{
					&onCall{stub: mockFn, index: 1},
					&onCall{stub: mockFn, index: 2},
					&onCall{stub: mockFn, index: 0},
				},
			}

			result := ca.OnThirdCall()

			o, ok := result.(*onCall)

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockFn, index: 2}))

		})
	})
})
