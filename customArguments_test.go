package mocka

import (
	"errors"
	"reflect"

	"github.com/MonsantoCo/mocka/match"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CustomArguments", func() {
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

	Describe("newCustomArguments", func() {
		It("returns with validation error when stub is nil", func() {
			ca := newCustomArguments(nil, []interface{}{})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(CustomArguments{
				argValidationError: &argumentValidationError{
					provided: []interface{}{},
				},
			}))
		})

		It("returns with validation error when stub is not a function type", func() {
			stub := &Stub{functionPtr: &struct{}{}}
			ca := newCustomArguments(stub, []interface{}{})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(CustomArguments{
				argValidationError: &argumentValidationError{
					provided: []interface{}{},
				},
			}))
		})

		It("returns with validation error when provided arguments != stubbed function arguments length", func() {
			ca := newCustomArguments(stub, []interface{}{})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(CustomArguments{
				argValidationError: &argumentValidationError{
					fnType:   stub.toType(),
					provided: []interface{}{},
				},
			}))
		})

		It("returns with validation error if the provided matcher is not supported for the argument kind", func() {
			anything := match.Anything()
			lengthOf10 := match.LengthOf(10)
			ca := newCustomArguments(stub, []interface{}{
				anything,
				lengthOf10,
			})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(CustomArguments{
				stub:        stub,
				argMatchers: nil,
				argValidationError: &argumentValidationError{
					fnType:   stub.toType(),
					provided: []interface{}{anything, lengthOf10},
				},
			}))
		})

		It("returns with validation error if the provided argument is not of the correct type", func() {
			ca := newCustomArguments(stub, []interface{}{"hi", "ope"})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(CustomArguments{
				stub:        stub,
				argMatchers: nil,
				argValidationError: &argumentValidationError{
					fnType:   stub.toType(),
					provided: []interface{}{"hi", "ope"},
				},
			}))
		})

		It("returns a valid custom arguments structs", func() {
			ca := newCustomArguments(stub, []interface{}{"hi", match.IntGreaterThan(10)})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(CustomArguments{
				stub:        stub,
				argMatchers: []match.SupportedKindsMatcher{match.Exactly("hi"), match.IntGreaterThan(10)},
			}))
		})

		It("returns with validation error if the provided nil does not match the correct type", func() {
			fn := func(msg string) error {
				return errors.New(msg)
			}
			stub = &Stub{
				originalFunc:  nil,
				functionPtr:   &fn,
				outParameters: []interface{}{42, nil},
				execFunc:      func([]interface{}) {},
			}

			ca := newCustomArguments(stub, []interface{}{nil})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(CustomArguments{
				stub:        stub,
				argMatchers: nil,
				argValidationError: &argumentValidationError{
					fnType:   stub.toType(),
					provided: []interface{}{nil},
				},
			}))
		})

		It("returns a nil matcher when called with nil", func() {
			fn := func(err error) error {
				return err
			}
			stub = &Stub{
				originalFunc:  nil,
				functionPtr:   &fn,
				outParameters: []interface{}{42, nil},
				execFunc:      func([]interface{}) {},
			}

			ca := newCustomArguments(stub, []interface{}{nil})

			Expect(ca).ToNot(BeNil())
			Expect(*ca).To(Equal(CustomArguments{
				stub:        stub,
				argMatchers: []match.SupportedKindsMatcher{match.Nil()},
			}))
		})

		Context("variadic function", func() {
			BeforeEach(func() {
				var variadicFn func(string, ...interface{}) (int, error)
				stub = &Stub{
					originalFunc:  nil,
					functionPtr:   &variadicFn,
					outParameters: []interface{}{42, nil},
					execFunc:      func([]interface{}) {},
				}
			})

			It("returns a CustomArguments struct with a nil matcher for omitted variadic arguments", func() {
				ca := newCustomArguments(stub, []interface{}{"hi"})

				Expect(ca).ToNot(BeNil())
				Expect(*ca).To(Equal(CustomArguments{
					stub:        stub,
					argMatchers: []match.SupportedKindsMatcher{match.Exactly("hi"), match.Nil()},
				}))
			})

			It("returns a CustomArguments struct with a sliceOf matcher for variadic arguments", func() {
				ca := newCustomArguments(stub, []interface{}{"hi", nil, "A", match.Anything()})

				Expect(ca).ToNot(BeNil())
				Expect(*ca).To(Equal(CustomArguments{
					stub: stub,
					argMatchers: []match.SupportedKindsMatcher{
						match.Exactly("hi"),
						match.SliceOf(match.Nil(), match.Exactly("A"), match.Anything())},
				}))
			})

			It("returns an argument validation error if matcher does not suppor the variadic type", func() {
				ca := newCustomArguments(stub, []interface{}{"hi", match.ElementsContaining("A")})

				Expect(ca).ToNot(BeNil())
				Expect(*ca).To(Equal(CustomArguments{
					stub:        stub,
					argMatchers: nil,
					argValidationError: &argumentValidationError{
						fnType:   stub.toType(),
						provided: []interface{}{"hi", match.ElementsContaining("A")},
					},
				}))
			})
		})
	})

	Describe("Return", func() {
		It("returns an error if the stub is nil", func() {
			ca := &CustomArguments{}

			err := ca.Return(42, "nil")

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("mocka: stub does not exist"))
		})

		It("returns an error if there is an argument validation error", func() {
			ca := &CustomArguments{
				stub: stub,
				argValidationError: &argumentValidationError{
					fnType:   stub.toType(),
					provided: []interface{}{10, 10},
				},
			}

			err := ca.Return(42, "nil")

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("mocka: expected arguments of type (string, int), but received (int, int)"))
		})

		It("returns an error if one or more of the return values are not valid", func() {
			ca := &CustomArguments{
				stub: stub,
				argMatchers: []match.SupportedKindsMatcher{
					match.Exactly(""), match.Exactly(42),
				},
			}

			err := ca.Return("", 42)

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("mocka: expected return values of type (int, error), but received (string, int)"))
		})

		It("assigns the OutParameters and returns nil if everything is valid", func() {
			ca := &CustomArguments{
				stub: stub,
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
		var ca *CustomArguments

		BeforeEach(func() {
			ca = &CustomArguments{
				stub: stub,
				onCalls: []*OnCall{
					&OnCall{stub: stub, index: 1},
					&OnCall{stub: stub, index: 2},
					&OnCall{stub: stub, index: 0},
				},
			}
		})

		It("returns a pointer to an onCall struct", func() {
			result := ca.OnCall(5)

			Expect(result).To(Equal(&OnCall{stub: stub, index: 5}))
		})

		It("appends the new onCall struct to the onCalls slice", func() {
			_ = ca.OnCall(5)

			Expect(ca.onCalls).To(HaveLen(4))
		})

		It("returns an existing onCall object if one exists for that index", func() {
			result := ca.OnCall(2)

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(result).To(Equal(&OnCall{stub: stub, index: 2}))
		})
	})

	Describe("OnFirstCall", func() {
		It("creates a new onCall with a 0 index", func() {
			ca := &CustomArguments{
				stub: stub,
				onCalls: []*OnCall{
					&OnCall{stub: stub, index: 1},
					&OnCall{stub: stub, index: 2},
				},
			}

			result := ca.OnFirstCall()

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(result).To(Equal(&OnCall{stub: stub, index: 0}))
		})

		It("returns the existing first call", func() {
			ca := &CustomArguments{
				stub: stub,
				onCalls: []*OnCall{
					&OnCall{stub: stub, index: 1},
					&OnCall{stub: stub, index: 2},
					&OnCall{stub: stub, index: 0},
				},
			}

			result := ca.OnFirstCall()

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(result).To(Equal(&OnCall{stub: stub, index: 0}))
		})
	})

	Describe("OnSecondCall", func() {
		It("creates a new onCall with a 1 index", func() {
			ca := &CustomArguments{
				stub: stub,
				onCalls: []*OnCall{
					&OnCall{stub: stub, index: 0},
					&OnCall{stub: stub, index: 2},
				},
			}

			result := ca.OnSecondCall()

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(result).To(Equal(&OnCall{stub: stub, index: 1}))
		})

		It("returns the existing second call", func() {
			ca := &CustomArguments{
				stub: stub,
				onCalls: []*OnCall{
					&OnCall{stub: stub, index: 1},
					&OnCall{stub: stub, index: 2},
					&OnCall{stub: stub, index: 0},
				},
			}

			result := ca.OnSecondCall()

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(result).To(Equal(&OnCall{stub: stub, index: 1}))
		})
	})

	Describe("OnThirdCall", func() {
		It("creates a new onCall with a 2 index", func() {
			ca := &CustomArguments{
				stub: stub,
				onCalls: []*OnCall{
					&OnCall{stub: stub, index: 0},
					&OnCall{stub: stub, index: 1},
				},
			}

			result := ca.OnThirdCall()

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(result).To(Equal(&OnCall{stub: stub, index: 2}))
		})

		It("returns the existing third call", func() {
			ca := &CustomArguments{
				stub: stub,
				onCalls: []*OnCall{
					&OnCall{stub: stub, index: 1},
					&OnCall{stub: stub, index: 2},
					&OnCall{stub: stub, index: 0},
				},
			}

			result := ca.OnThirdCall()

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(result).To(Equal(&OnCall{stub: stub, index: 2}))
		})
	})

	Describe("isMatch", func() {
		It("returns false if any matcher panics", func() {
			ca := newCustomArguments(stub, []interface{}{&panicMatcher{}, match.IntGreaterThan(10)})

			Expect(ca.isMatch([]interface{}{"hi", 11})).To(BeFalse())
		})

		It("returns false if any matcher returns false", func() {
			ca := newCustomArguments(stub, []interface{}{"hi", match.IntGreaterThan(10)})

			Expect(ca.isMatch([]interface{}{"hi", 5})).To(BeFalse())
		})

		It("returns true if all matchers return true", func() {
			ca := newCustomArguments(stub, []interface{}{"hi", match.IntGreaterThan(10)})

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
