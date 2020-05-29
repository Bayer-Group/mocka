package mocka

import (
	"errors"
	"reflect"

	"github.com/MonsantoCo/mocka/match"
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

	Describe("newCustomArguments", func() {
		It("returns with validation error when stub is nil", func() {
			ca := newCustomArguments(nil, []interface{}{})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(customArguments{
				argValidationError: &argumentValidationError{
					provided: []interface{}{},
				},
			}))
		})

		It("returns with validation error when stub is not a function type", func() {
			stub := &mockFunction{functionPtr: &struct{}{}}
			ca := newCustomArguments(stub, []interface{}{})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(customArguments{
				argValidationError: &argumentValidationError{
					provided: []interface{}{},
				},
			}))
		})

		It("returns with validation error when provided arguments != stubbed function arguments length", func() {
			ca := newCustomArguments(mockFn, []interface{}{})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(customArguments{
				argValidationError: &argumentValidationError{
					fnType:   mockFn.toType(),
					provided: []interface{}{},
				},
			}))
		})

		It("returns with validation error if the provided matcher is not supported for the argument kind", func() {
			anything := match.Anything()
			lengthOf10 := match.LengthOf(10)
			ca := newCustomArguments(mockFn, []interface{}{
				anything,
				lengthOf10,
			})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(customArguments{
				stub:        mockFn,
				argMatchers: nil,
				argValidationError: &argumentValidationError{
					fnType:   mockFn.toType(),
					provided: []interface{}{anything, lengthOf10},
				},
			}))
		})

		It("returns with validation error if the provided argument is not of the correct type", func() {
			ca := newCustomArguments(mockFn, []interface{}{"hi", "ope"})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(customArguments{
				stub:        mockFn,
				argMatchers: nil,
				argValidationError: &argumentValidationError{
					fnType:   mockFn.toType(),
					provided: []interface{}{"hi", "ope"},
				},
			}))
		})

		It("returns a valid custom arguments structs", func() {
			ca := newCustomArguments(mockFn, []interface{}{"hi", match.IntGreaterThan(10)})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(customArguments{
				stub:        mockFn,
				argMatchers: []match.SupportedKindsMatcher{match.Exactly("hi"), match.IntGreaterThan(10)},
			}))
		})

		It("returns with validation error if the provided nil does not match the correct type", func() {
			fn := func(msg string) error {
				return errors.New(msg)
			}
			mockFn = &mockFunction{
				originalFunc:  nil,
				functionPtr:   &fn,
				outParameters: []interface{}{42, nil},
				execFunc:      func([]interface{}) {},
			}

			ca := newCustomArguments(mockFn, []interface{}{nil})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(customArguments{
				stub:        mockFn,
				argMatchers: nil,
				argValidationError: &argumentValidationError{
					fnType:   mockFn.toType(),
					provided: []interface{}{nil},
				},
			}))
		})

		It("returns a nil matcher when called with nil", func() {
			fn := func(err error) error {
				return err
			}
			mockFn = &mockFunction{
				originalFunc:  nil,
				functionPtr:   &fn,
				outParameters: []interface{}{42, nil},
				execFunc:      func([]interface{}) {},
			}

			ca := newCustomArguments(mockFn, []interface{}{nil})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(customArguments{
				stub:        mockFn,
				argMatchers: []match.SupportedKindsMatcher{match.Nil()},
			}))
		})
	})

	Describe("Return", func() {
		It("returns an error if the stub is nil", func() {
			ca := &customArguments{}

			err := ca.Return(42, "nil")

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("mocka: stub does not exist"))
		})

		It("returns an error if there is an argument validation error", func() {
			ca := &customArguments{
				stub: mockFn,
				argValidationError: &argumentValidationError{
					fnType:   mockFn.toType(),
					provided: []interface{}{10, 10},
				},
			}

			err := ca.Return(42, "nil")

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("mocka: expected arguments of type (string, int), but received (int, int)"))
		})

		It("returns an error if one or more of the return values are not valid", func() {
			ca := &customArguments{
				stub: mockFn,
				argMatchers: []match.SupportedKindsMatcher{
					match.Exactly(""), match.Exactly(42),
				},
			}

			err := ca.Return("", 42)

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("mocka: expected return values of type (int, error), but received (string, int)"))
		})

		It("assigns the OutParameters and returns nil if everything is valid", func() {
			ca := &customArguments{
				stub: mockFn,
				argMatchers: []match.SupportedKindsMatcher{
					match.Exactly(""), match.Exactly(42),
				},
			}

			err := ca.Return(42, nil)

			Expect(err).ShouldNot(HaveOccurred())
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

	Describe("match", func() {
		It("returns false if any matcher panics", func() {
			ca := newCustomArguments(mockFn, []interface{}{&panicMatcher{}, match.IntGreaterThan(10)})

			Expect(ca.isMatch([]interface{}{"hi", 11})).To(BeFalse())
		})

		It("returns false if any matcher returns false", func() {
			ca := newCustomArguments(mockFn, []interface{}{"hi", match.IntGreaterThan(10)})

			Expect(ca.isMatch([]interface{}{"hi", 5})).To(BeFalse())
		})

		It("returns true if all matchers return true", func() {
			ca := newCustomArguments(mockFn, []interface{}{"hi", match.IntGreaterThan(10)})

			Expect(ca.isMatch([]interface{}{"hi", 15})).To(BeTrue())
		})
	})
})

type panicMatcher struct {
}

func (*panicMatcher) SupportedKinds() map[reflect.Kind]struct{} {
	return map[reflect.Kind]struct{}{
		reflect.Int:    {},
		reflect.String: {},
	}
}

func (*panicMatcher) Match(_ interface{}) bool {
	panic("i am the panic matcher")
}
