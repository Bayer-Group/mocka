package mocka

import (
	"errors"
	"reflect"

	"github.com/MonsantoCo/mocka/v2/match"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CustomArguments", func() {
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

	Describe("newCustomArguments", func() {
		It("reports an error when provided arguments != stubbed function arguments length", func() {
			stub.testReporter = failTestReporter

			_ = newCustomArguments(stub, []interface{}{})

			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected arguments of type (string, int), but received ()",
			}))
		})

		It("reports an error if the provided matcher is not supported for the argument kind", func() {
			stub.testReporter = failTestReporter
			anything := match.Anything()
			lengthOf10 := match.LengthOf(10)

			_ = newCustomArguments(stub, []interface{}{
				anything,
				lengthOf10,
			})

			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected arguments of type (string, int), but received (*anything, *lengthOf)",
			}))
		})

		It("reports an error if the provided argument is not of the correct type", func() {
			stub.testReporter = failTestReporter

			_ = newCustomArguments(stub, []interface{}{"hi", "ope"})

			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected arguments of type (string, int), but received (string, string)",
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

		It("reports an error if the provided nil does not match the correct type", func() {
			fn := func(msg string) error {
				return errors.New(msg)
			}
			stub = &Stub{
				testReporter:  failTestReporter,
				originalFunc:  nil,
				functionPtr:   &fn,
				outParameters: []interface{}{42, nil},
				execFunc:      func([]interface{}) {},
			}

			_ = newCustomArguments(stub, []interface{}{nil})

			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected arguments of type (string), but received (<nil>)",
			}))
		})

		It("returns a nil matcher when called with nil", func() {
			fn := func(err error) error {
				return err
			}
			stub = &Stub{
				testReporter:  GinkgoT(),
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
					testReporter:  GinkgoT(),
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

			It("reports an error if matcher does not suppor the variadic type", func() {
				stub.testReporter = failTestReporter

				_ = newCustomArguments(stub, []interface{}{"hi", match.ElementsContaining("A")})

				Expect(failTestReporter.messages).To(Equal([]string{
					"mocka: expected arguments of type (string, ...), but received (string, *elementsContaining)",
				}))
			})
		})
	})

	Describe("Return", func() {
		It("returns an error if one or more of the return values are not valid", func() {
			stub.testReporter = failTestReporter
			ca := &CustomArguments{
				stub: stub,
				argMatchers: []match.SupportedKindsMatcher{
					match.Exactly(""), match.Exactly(42),
				},
			}

			ca.Return("", 42)

			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected return values of type (int, error), but received (string, int)",
			}))
		})

		It("assigns the OutParameters and returns nil if everything is valid", func() {
			ca := &CustomArguments{
				stub: stub,
				argMatchers: []match.SupportedKindsMatcher{
					match.Exactly(""), match.Exactly(42),
				},
			}

			ca.Return(42, nil)

			Expect(ca.out).To(Equal([]interface{}{42, nil}))
		})
	})

	Describe("OnCall", func() {
		var ca *CustomArguments

		BeforeEach(func() {
			ca = &CustomArguments{
				stub: stub,
				onCalls: []*OnCall{
					{stub: stub, index: 1},
					{stub: stub, index: 2},
					{stub: stub, index: 0},
				},
			}
		})

		It("returns a pointer to an onCall struct", func() {
			result := ca.OnCall(5)

			Expect(*result).To(Equal(OnCall{stub: stub, index: 5}))
		})

		It("appends the new onCall struct to the onCalls slice", func() {
			_ = ca.OnCall(5)

			Expect(ca.onCalls).To(HaveLen(4))
		})

		It("returns an existing onCall object if one exists for that index", func() {
			result := ca.OnCall(2)

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(*result).To(Equal(OnCall{stub: stub, index: 2}))
		})
	})

	Describe("OnFirstCall", func() {
		It("creates a new onCall with a 0 index", func() {
			ca := &CustomArguments{
				stub: stub,
				onCalls: []*OnCall{
					{stub: stub, index: 1},
					{stub: stub, index: 2},
				},
			}

			result := ca.OnFirstCall()

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(*result).To(Equal(OnCall{stub: stub, index: 0}))
		})

		It("returns the existing first call", func() {
			ca := &CustomArguments{
				stub: stub,
				onCalls: []*OnCall{
					{stub: stub, index: 1},
					{stub: stub, index: 2},
					{stub: stub, index: 0},
				},
			}

			result := ca.OnFirstCall()

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(*result).To(Equal(OnCall{stub: stub, index: 0}))
		})
	})

	Describe("OnSecondCall", func() {
		It("creates a new onCall with a 1 index", func() {
			ca := &CustomArguments{
				stub: stub,
				onCalls: []*OnCall{
					{stub: stub, index: 0},
					{stub: stub, index: 2},
				},
			}

			result := ca.OnSecondCall()

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(*result).To(Equal(OnCall{stub: stub, index: 1}))
		})

		It("returns the existing second call", func() {
			ca := &CustomArguments{
				stub: stub,
				onCalls: []*OnCall{
					{stub: stub, index: 1},
					{stub: stub, index: 2},
					{stub: stub, index: 0},
				},
			}

			result := ca.OnSecondCall()

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(*result).To(Equal(OnCall{stub: stub, index: 1}))
		})
	})

	Describe("OnThirdCall", func() {
		It("creates a new onCall with a 2 index", func() {
			ca := &CustomArguments{
				stub: stub,
				onCalls: []*OnCall{
					{stub: stub, index: 0},
					{stub: stub, index: 1},
				},
			}

			result := ca.OnThirdCall()

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(*result).To(Equal(OnCall{stub: stub, index: 2}))
		})

		It("returns the existing third call", func() {
			ca := &CustomArguments{
				stub: stub,
				onCalls: []*OnCall{
					{stub: stub, index: 1},
					{stub: stub, index: 2},
					{stub: stub, index: 0},
				},
			}

			result := ca.OnThirdCall()

			Expect(ca.onCalls).To(HaveLen(3))
			Expect(*result).To(Equal(OnCall{stub: stub, index: 2}))
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
