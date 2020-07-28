package mocka

import (
	"errors"
	"reflect"

	"github.com/MonsantoCo/mocka/v2/match"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("stub", func() {
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

		err := cloneValue(&fn, &stub.originalFunc)
		Expect(err).To(Succeed())
		failTestReporter = &mockTestReporter{}
	})

	AfterEach(func() {
		// clean up pointers and slices to prevent memory leaks
		stub.originalFunc = nil
		stub.functionPtr = nil
		stub.outParameters = nil
		stub.calls = nil
		stub.customArgs = nil
		stub = nil
	})

	Describe("newStub", func() {
		var callCount int

		BeforeEach(func() {
			callCount = 0
			fn = func(str string, num int) (int, error) {
				callCount++
				return len(str) + num, nil
			}
		})

		It("reports an error if passed a nil as the function pointer", func() {
			stub := newStub(failTestReporter, nil, nil)

			Expect(stub).To(BeNil())
			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected the second argument to be a pointer to a function, but received a nil",
			}))
		})

		It("reports an error if a non-pointer value is passed as the function pointer", func() {
			stub := newStub(failTestReporter, 42, nil)

			Expect(stub).To(BeNil())
			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected the second argument to be a pointer to a function, but received a int",
			}))
		})

		It("reports an error if a non-function value is passed as the function pointer", func() {
			num := 42
			stub := newStub(failTestReporter, &num, nil)

			Expect(stub).To(BeNil())
			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected the second argument to be a pointer to a function, but received a pointer to a int",
			}))
		})

		It("reports an error supplied out parameters are not of the same type", func() {
			stub := newStub(failTestReporter, &fn, []interface{}{"42", nil})

			Expect(stub).To(BeNil())
			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected return values of type (int, error), but received (string, <nil>)",
			}))
		})

		It("reports an error if cloneValue returns an error", func() {
			_cloneValue = func(interface{}, interface{}) error {
				return errors.New("Ope")
			}
			defer func() {
				_cloneValue = cloneValue
			}()

			stub := newStub(failTestReporter, &fn, []interface{}{42, nil})

			Expect(stub).To(BeNil())
			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: could not clone function pointer to new memory address: Ope",
			}))
		})

		It("returns a Stub with a reference to the original function", func() {
			stub := newStub(GinkgoT(), &fn, []interface{}{42, nil})

			Expect(stub).ToNot(BeNil())
			Expect(stub.originalFunc).ToNot(BeNil())

			_, _ = stub.originalFunc.(func(str string, num int) (int, error))("", 0)

			Expect(callCount).To(Equal(1))
		})

		It("returns a Stub with properties initialized with zero values", func() {
			stub := newStub(GinkgoT(), &fn, []interface{}{42, nil})

			Expect(stub).ToNot(BeNil())
			Expect(stub.calls).To(BeNil())
			Expect(stub.customArgs).To(BeNil())
		})

		It("returns a Stub with outParameters as supplied", func() {
			stub := newStub(GinkgoT(), &fn, []interface{}{42, nil})

			Expect(stub).ToNot(BeNil())
			Expect(stub.outParameters).To(Equal([]interface{}{42, nil}))
		})
	})

	Describe("getReturnValues", func() {
		It("returns the Stub.OutParameters if no customArgs or onCalls exist", func() {
			args := []interface{}{"Hello", 42}

			result, maybeCustomArguments := stub.getReturnValues(args, reflect.TypeOf(fn))

			Expect(result).To(Equal([]interface{}{42, nil}))
			Expect(maybeCustomArguments).To(BeNil())
		})

		It("returns the Stub.OutParameters if the customArgs are nil", func() {
			args := []interface{}{"Hello", 42}
			stub.customArgs = append(stub.customArgs, nil, nil, nil)

			result, maybeCustomArguments := stub.getReturnValues(args, reflect.TypeOf(fn))

			Expect(result).To(Equal([]interface{}{42, nil}))
			Expect(maybeCustomArguments).To(BeNil())
		})

		It("returns the Stub.OutParameters if the argument do not equal any customArgs", func() {
			args := []interface{}{"Hello", 42}
			stub.customArgs = append(
				stub.customArgs,
				&CustomArguments{
					stub:        stub,
					argMatchers: []match.SupportedKindsMatcher{match.Exactly("apple"), match.IntGreaterThanOrEqualTo(0)},
					out:         []interface{}{0, errors.New("I am not an apple")},
				},
				&CustomArguments{
					stub:        stub,
					argMatchers: []match.SupportedKindsMatcher{match.Exactly("B"), match.Exactly(63)},
					out:         []interface{}{98, nil},
				})

			result, maybeCustomArguments := stub.getReturnValues(args, reflect.TypeOf(fn))

			Expect(result).To(Equal([]interface{}{42, nil}))
			Expect(maybeCustomArguments).To(BeNil())
		})

		It("returns the out parameters if the customArgs if the arguments match", func() {
			args := []interface{}{"Hello", 42}
			expected := &CustomArguments{
				stub:        stub,
				argMatchers: []match.SupportedKindsMatcher{match.StringSuffix("ello"), match.Exactly(42)},
				out:         []interface{}{22, errors.New("I am an error")},
			}
			stub.customArgs = append(
				stub.customArgs,
				&CustomArguments{
					stub:        stub,
					argMatchers: []match.SupportedKindsMatcher{match.StringPrefix("app"), match.Exactly(0)},
					out:         []interface{}{0, errors.New("I am not an apple")},
				},
				&CustomArguments{
					stub:        stub,
					argMatchers: []match.SupportedKindsMatcher{match.Exactly("B"), match.IntGreaterThan(60)},
					out:         []interface{}{98, nil},
				},
				expected,
			)

			result, maybeCustomArguments := stub.getReturnValues(args, reflect.TypeOf(fn))

			Expect(result).To(Equal([]interface{}{22, errors.New("I am an error")}))
			Expect(maybeCustomArguments).To(Equal(expected))
		})

		It("returns the the out parameters for the specific call index", func() {
			args := []interface{}{"apples", 42}
			stub.customArgs = append(
				stub.customArgs,
				&CustomArguments{
					stub:        stub,
					argMatchers: []match.SupportedKindsMatcher{match.Exactly("apple"), match.Exactly(0)},
					out:         []interface{}{0, errors.New("I am not an apple")},
					callCount:   2,
					onCalls: []*OnCall{
						{
							stub:  stub,
							index: 2,
							out:   []interface{}{23, errors.New("I am the third not an apple")},
						},
					},
				},
			)
			stub.calls = []Call{{}, {}, {}}
			stub.onCalls = append(stub.onCalls, &OnCall{
				stub:  stub,
				index: 3,
				out:   []interface{}{22, errors.New("I am the first error")},
			})

			result, maybeCustomArguments := stub.getReturnValues(args, reflect.TypeOf(fn))

			Expect(result).To(Equal([]interface{}{22, errors.New("I am the first error")}))
			Expect(maybeCustomArguments).To(BeNil())
		})

		It("returns the the out parameters for the specific call index on a custom argument", func() {
			args := []interface{}{"apple", 0}
			expected := &CustomArguments{
				stub:        stub,
				argMatchers: []match.SupportedKindsMatcher{match.Exactly("apple"), match.Exactly(0)},
				out:         []interface{}{0, errors.New("I am not an apple")},
				callCount:   2,
				onCalls: []*OnCall{
					{
						stub:  stub,
						index: 2,
						out:   []interface{}{23, errors.New("I am the third not an apple")},
					},
				},
			}
			stub.customArgs = append(stub.customArgs, expected)
			stub.calls = []Call{{}, {}, {}}
			stub.onCalls = append(stub.onCalls, &OnCall{
				stub:  stub,
				index: 0,
				out:   []interface{}{22, errors.New("I am the first error")},
			})

			result, maybeCustomArguments := stub.getReturnValues(args, reflect.TypeOf(fn))

			Expect(result).To(Equal([]interface{}{23, errors.New("I am the third not an apple")}))
			Expect(maybeCustomArguments).To(Equal(expected))
		})
	})

	Describe("toType", func() {
		It("returns the tuype of the mocked function", func() {
			fnValue := reflect.ValueOf(&fn).Elem()

			result := stub.toType()

			Expect(result).To(Equal(fnValue.Type()))
		})
	})

	Describe("implementation", func() {
		BeforeEach(func() {
			stub.outParameters = []interface{}{42, nil}
			stub.calls = []Call{}
			stub.customArgs = []*CustomArguments{
				{
					stub:        stub,
					argMatchers: []match.SupportedKindsMatcher{match.Exactly("custom"), match.Exactly(0)},
					out:         []interface{}{0, errors.New("Ope")},
				},
				{
					stub:        stub,
					argMatchers: []match.SupportedKindsMatcher{match.StringPrefix("custom-"), match.Exactly(0)},
					out:         nil,
				},
			}
		})

		It("calls the exec function with the arguments as interfaces", func() {
			var argsProvided []interface{}
			args := []reflect.Value{reflect.ValueOf("Hello"), reflect.ValueOf(42)}
			expected := []interface{}{"Hello", 42}

			stub.execFunc = func(a []interface{}) {
				argsProvided = a
			}

			_ = stub.implementation(args)

			Expect(argsProvided).To(Equal(expected))
		})

		It("appends the call meta data for this call into the calls slice", func() {
			args := []reflect.Value{reflect.ValueOf("Hello"), reflect.ValueOf(42)}

			Expect(stub.calls).To(HaveLen(0))

			_ = stub.implementation(args)

			Expect(stub.calls).To(HaveLen(1))

			ca := stub.calls[0]

			Expect(ca.args).To(Equal([]interface{}{"Hello", 42}))
			Expect(ca.out).To(Equal([]interface{}{42, nil}))
		})

		It("returns the out parameters as reflection values", func() {
			args := []reflect.Value{reflect.ValueOf("Hello"), reflect.ValueOf(42)}

			outValues := stub.implementation(args)
			outInterfaces := mapToInterfaces(outValues)

			Expect(outInterfaces).To(Equal([]interface{}{42, nil}))
		})

		It("uses out parameters for custom arguments if applicable", func() {
			args := []reflect.Value{reflect.ValueOf("custom"), reflect.ValueOf(0)}

			outValues := stub.implementation(args)
			outInterfaces := mapToInterfaces(outValues)

			Expect(outInterfaces).To(Equal([]interface{}{0, errors.New("Ope")}))
		})

		It("doesn't use custom arguments if the return values have not been supplied", func() {
			args := []reflect.Value{reflect.ValueOf("custom-missing"), reflect.ValueOf(0)}

			outValues := stub.implementation(args)
			outInterfaces := mapToInterfaces(outValues)

			Expect(outInterfaces).To(Equal([]interface{}{42, nil}))
		})

		Context("variadic function", func() {
			BeforeEach(func() {
				fn := func(str string, opts ...string) (int, error) {
					return len(str) + len(opts), nil
				}

				stub.functionPtr = &fn
				Expect(cloneValue(&fn, &stub.originalFunc)).To(Succeed())
				stub.outParameters = []interface{}{42, nil}
				stub.calls = []Call{}
			})

			It("appends the call meta data omitting the missing variadic arguments", func() {
				args := []reflect.Value{reflect.ValueOf("Hello")}

				Expect(stub.calls).To(HaveLen(0))

				_ = stub.implementation(args)

				Expect(stub.calls).To(HaveLen(1))

				ca := stub.calls[0]

				Expect(ca.args).To(Equal([]interface{}{"Hello"}))
			})

			It("appends the call meta data spreading the variadic arguments", func() {
				args := []reflect.Value{reflect.ValueOf("Hello"), reflect.ValueOf("A"), reflect.ValueOf("B")}

				Expect(stub.calls).To(HaveLen(0))

				_ = stub.implementation(args)

				Expect(stub.calls).To(HaveLen(1))

				ca := stub.calls[0]

				Expect(ca.args).To(Equal([]interface{}{"Hello", "A", "B"}))
			})
		})
	})

	Describe("Return", func() {
		It("reports an error if the out parameters are not valid", func() {
			stub.testReporter = failTestReporter

			stub.Return(42, 42)

			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected return values of type (int, error), but received (int, int)",
			}))
		})

		It("replaces the out parameters if they are valid", func() {
			stub.Return(22, errors.New("I am new"))

			Expect(stub.outParameters).To(Equal([]interface{}{22, errors.New("I am new")}))
		})
	})

	Describe("WithArgs", func() {
		It("returns existing custom arguments if it matching arguments", func() {
			ca := &CustomArguments{
				stub:        stub,
				argMatchers: []match.SupportedKindsMatcher{match.Exactly("apple"), match.Exactly(0)},
				out:         []interface{}{0, errors.New("I am not an apple")},
			}
			stub.customArgs = append(stub.customArgs, ca)

			withArgs := stub.WithArgs("apple", 0)

			Expect(withArgs).ToNot(BeNil())
			Expect(withArgs).To(Equal(ca))
		})

		It("creates and returns new custom arguments if one does not exist", func() {
			Expect(stub.customArgs).To(HaveLen(0))

			withArgs := stub.WithArgs("apple", 0)

			Expect(withArgs).ToNot(BeNil())
			Expect(stub.customArgs).To(HaveLen(1))

			Expect(withArgs.stub).To(Equal(stub))
			Expect(withArgs.argMatchers).To(Equal([]match.SupportedKindsMatcher{match.Exactly("apple"), match.Exactly(0)}))
			Expect(withArgs.out).To(BeNil())
		})
	})

	Describe("CallCount", func() {
		It("returns the number of times the stub was called", func() {
			Expect(stub.CallCount()).To(Equal(0))

			stub.calls = append(stub.calls, Call{})

			Expect(stub.CallCount()).To(Equal(1))

			stub.calls = append(stub.calls, Call{})

			Expect(stub.CallCount()).To(Equal(2))
		})
	})

	Describe("GetCalls", func() {
		It("returns all the call meta data made to the stub", func() {
			stub.calls = append(stub.calls,
				Call{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				Call{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				Call{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			)

			result := stub.GetCalls()

			Expect(result).To(Equal([]Call{
				{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}))
		})
	})

	Describe("GetCall", func() {
		BeforeEach(func() {
			stub.calls = []Call{
				{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}
		})

		It("panics if the call index is less than 0", func() {
			stub.testReporter = failTestReporter

			_ = stub.GetCall(-1)

			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: attempted to get Call for invocation -1, when the function has only been called 3 times",
			}))
		})

		It("panics if the call index is greater than the number of calls", func() {
			stub.testReporter = failTestReporter

			_ = stub.GetCall(5)

			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: attempted to get Call for invocation 5, when the function has only been called 3 times",
			}))
		})

		It("returns the call meta data for the specified call index", func() {
			result := stub.GetCall(1)

			Expect(result).To(Equal(Call{
				args: []interface{}{"sam", 22},
				out:  []interface{}{42, nil},
			}))
		})
	})

	Describe("GetFirstCall", func() {
		It("panics if the stub has not been called once", func() {
			stub.testReporter = failTestReporter

			_ = stub.GetFirstCall()

			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: attempted to get Call for invocation 0, when the function has only been called 0 times",
			}))
		})

		It("returns the call meta data for the first call", func() {
			stub.calls = []Call{
				{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}

			result := stub.GetFirstCall()

			Expect(result).To(Equal(Call{
				args: []interface{}{"hello", 42},
				out:  []interface{}{22, nil},
			}))
		})
	})

	Describe("GetSecondCall", func() {
		It("panics if the stub has not been called twice", func() {
			stub.testReporter = failTestReporter

			_ = stub.GetSecondCall()

			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: attempted to get Call for invocation 1, when the function has only been called 0 times",
			}))
		})

		It("returns the call meta data for the second call", func() {
			stub.calls = []Call{
				{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}

			result := stub.GetSecondCall()

			Expect(result).To(Equal(Call{
				args: []interface{}{"sam", 22},
				out:  []interface{}{42, nil},
			}))
		})
	})

	Describe("GetThirdCall", func() {
		It("panics if the stub has not been called three times", func() {
			stub.testReporter = failTestReporter

			_ = stub.GetThirdCall()

			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: attempted to get Call for invocation 2, when the function has only been called 0 times",
			}))
		})

		It("returns the call meta data for the third call", func() {
			stub.calls = []Call{
				{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}

			result := stub.GetThirdCall()

			Expect(result).To(Equal(Call{
				args: []interface{}{"rob", 12},
				out:  []interface{}{0, errors.New("ope")},
			}))
		})
	})

	Describe("CalledOnce", func() {
		It("returns true is the stub has been called at least once", func() {
			stub.calls = []Call{
				{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}

			result := stub.CalledOnce()

			Expect(result).To(BeTrue())
		})

		It("returns false if the stub has not been called at least once", func() {
			result := stub.CalledOnce()

			Expect(result).To(BeFalse())
		})
	})

	Describe("CalledTwice", func() {
		It("returns true is the stub has been called at least twice", func() {
			stub.calls = []Call{
				{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}

			result := stub.CalledTwice()

			Expect(result).To(BeTrue())
		})

		It("returns false if the stub has not been called at least twice", func() {
			result := stub.CalledTwice()

			Expect(result).To(BeFalse())
		})
	})

	Describe("CalledThrice", func() {
		It("returns true is the stub has been called at least thrice", func() {
			stub.calls = []Call{
				{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}

			result := stub.CalledThrice()

			Expect(result).To(BeTrue())
		})

		It("returns false if the stub has not been called at least thrice", func() {
			result := stub.CalledThrice()

			Expect(result).To(BeFalse())
		})
	})

	Describe("OnCall", func() {
		BeforeEach(func() {
			stub.onCalls = []*OnCall{
				{stub: stub, index: 1},
				{stub: stub, index: 2},
				{stub: stub, index: 0},
			}
		})

		It("returns a pointer to an onCall struct", func() {
			result := stub.OnCall(5)

			Expect(*result).To(Equal(OnCall{stub: stub, index: 5}))
		})

		It("appends the new onCall struct to the onCalls slice", func() {
			_ = stub.OnCall(5)

			Expect(stub.onCalls).To(HaveLen(4))
		})

		It("returns an existing onCall object if one exists for that index", func() {
			result := stub.OnCall(2)

			Expect(stub.onCalls).To(HaveLen(3))
			Expect(*result).To(Equal(OnCall{stub: stub, index: 2}))
		})
	})

	Describe("OnFirstCall", func() {
		It("creates a new onCall with a 0 index", func() {
			stub.onCalls = []*OnCall{
				{stub: stub, index: 1},
				{stub: stub, index: 2},
			}

			result := stub.OnFirstCall()

			Expect(stub.onCalls).To(HaveLen(3))
			Expect(*result).To(Equal(OnCall{stub: stub, index: 0}))
		})

		It("returns the existing first call", func() {
			stub.onCalls = []*OnCall{
				{stub: stub, index: 1},
				{stub: stub, index: 2},
				{stub: stub, index: 0},
			}

			result := stub.OnFirstCall()

			Expect(stub.onCalls).To(HaveLen(3))
			Expect(*result).To(Equal(OnCall{stub: stub, index: 0}))
		})
	})

	Describe("OnSecondCall", func() {
		It("creates a new onCall with a 1 index", func() {
			stub.onCalls = []*OnCall{
				{stub: stub, index: 0},
				{stub: stub, index: 2},
			}

			result := stub.OnSecondCall()

			Expect(stub.onCalls).To(HaveLen(3))
			Expect(*result).To(Equal(OnCall{stub: stub, index: 1}))
		})

		It("returns the existing second call", func() {
			stub.onCalls = []*OnCall{
				{stub: stub, index: 1},
				{stub: stub, index: 2},
				{stub: stub, index: 0},
			}

			result := stub.OnSecondCall()

			Expect(stub.onCalls).To(HaveLen(3))
			Expect(*result).To(Equal(OnCall{stub: stub, index: 1}))
		})
	})

	Describe("OnThirdCall", func() {
		It("creates a new onCall with a 2 index", func() {
			stub.onCalls = []*OnCall{
				{stub: stub, index: 0},
				{stub: stub, index: 1},
			}

			result := stub.OnThirdCall()

			Expect(stub.onCalls).To(HaveLen(3))
			Expect(*result).To(Equal(OnCall{stub: stub, index: 2}))
		})

		It("returns the existing third call", func() {
			stub.onCalls = []*OnCall{
				{stub: stub, index: 1},
				{stub: stub, index: 2},
				{stub: stub, index: 0},
			}

			result := stub.OnThirdCall()

			Expect(stub.onCalls).To(HaveLen(3))
			Expect(*result).To(Equal(OnCall{stub: stub, index: 2}))
		})
	})

	Describe("Restore", func() {
		It("overrides the function pointer with the the original pointer", func() {
			stub.originalFunc = func(str string, num int) (int, error) {
				return 42, errors.New("Ope")
			}

			castfn, ok := stub.functionPtr.(*func(str string, num int) (int, error))
			Expect(ok).To(BeTrue())

			n, err := (*castfn)("hello", 2)
			Expect(n).To(Equal(7))
			Expect(err).To(BeNil())

			stub.Restore()

			castfn, ok = stub.functionPtr.(*func(str string, num int) (int, error))
			Expect(ok).To(BeTrue())

			n, err = (*castfn)("hello", 2)
			Expect(n).To(Equal(42))
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("Ope"))
		})
	})

	Describe("ExecOnCall", func() {
		It("assigns the exec function to the new function provided", func() {
			called := false
			var newFn = func(vals []interface{}) {
				called = true
			}

			stub.ExecOnCall(newFn)

			Expect(called).To(BeFalse())

			stub.execFunc([]interface{}{})

			Expect(called).To(BeTrue())
		})
	})

	Describe("getHighestPriority", func() {
		const numArguments = 2

		var (
			customArgs []*CustomArguments
			matcher1   *CustomArguments
			matcher2   *CustomArguments
		)

		BeforeEach(func() {
			matcher1 = &CustomArguments{
				stub:        stub,
				argMatchers: []match.SupportedKindsMatcher{match.Exactly("custom-"), match.Exactly(0)},
				out:         []interface{}{0, errors.New("Ope")},
			}
			matcher2 = &CustomArguments{
				stub:        stub,
				argMatchers: []match.SupportedKindsMatcher{match.StringPrefix("custom-"), match.Exactly(0)},
				out:         nil,
			}
			customArgs = []*CustomArguments{matcher1, matcher2}
		})

		It("returns nil when no custom arguments are provided", func() {
			actual := getHighestPriority(nil, numArguments)

			Expect(actual).To(BeNil())
		})

		It("return the only custom argument when provided one custom argument", func() {
			actual := getHighestPriority([]*CustomArguments{matcher1}, numArguments)

			Expect(actual).To(Equal(matcher1))
		})

		It("returns the matcher with the highest priority", func() {
			actual := getHighestPriority(customArgs, numArguments)

			Expect(actual).To(Equal(matcher1))
		})

		It("returns the first matcher if multiple matchers have the same priority", func() {
			customArgs = []*CustomArguments{
				matcher2,
				{
					stub:        stub,
					argMatchers: []match.SupportedKindsMatcher{match.StringPrefix("custom"), match.Exactly(0)},
					out:         nil,
				},
			}
			actual := getHighestPriority(customArgs, numArguments)

			Expect(actual).To(Equal(matcher2))
		})
	})

	Describe("getPossible", func() {
		var (
			customArgs []*CustomArguments
			matcher1   *CustomArguments
			matcher2   *CustomArguments
		)

		BeforeEach(func() {
			matcher1 = &CustomArguments{
				stub:        stub,
				argMatchers: []match.SupportedKindsMatcher{match.Exactly("custom-"), match.Exactly(0)},
				out:         []interface{}{0, errors.New("Ope")},
			}
			matcher2 = &CustomArguments{
				stub:        stub,
				argMatchers: []match.SupportedKindsMatcher{match.StringPrefix("custom-"), match.Exactly(0)},
				out:         nil,
			}
			customArgs = []*CustomArguments{matcher1, matcher2}
		})

		It("returns an empty slice when no matches are found", func() {
			actual := getPossible(customArgs, []interface{}{"screams", 0})

			Expect(actual).To(HaveLen(0))
		})

		It("returns all possible matches", func() {
			actual := getPossible(customArgs, []interface{}{"custom-", 0})

			Expect(actual).To(HaveLen(2))
			Expect(actual).To(ContainElement(matcher1))
			Expect(actual).To(ContainElement(matcher2))
		})

		It("returns a single match when only one is possible", func() {
			actual := getPossible(customArgs, []interface{}{"custom-1", 0})

			Expect(actual).To(HaveLen(1))
			Expect(actual).To(ContainElement(matcher2))
		})
	})
})
