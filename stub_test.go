package mocka

import (
	"errors"
	"reflect"

	"github.com/MonsantoCo/mocka/match"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("stub", func() {
	var (
		fn     func(string, int) (int, error)
		mockfn *mockFunction
	)

	BeforeEach(func() {
		fn = func(str string, num int) (int, error) {
			return len(str) + num, nil
		}

		mockfn = &mockFunction{
			originalFunc:  nil,
			functionPtr:   &fn,
			outParameters: []interface{}{42, nil},
			execFunc:      func([]interface{}) {},
		}

		err := cloneValue(&fn, &mockfn.originalFunc)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		// clean up pointers and slices to prevent memory leaks
		mockfn.originalFunc = nil
		mockfn.functionPtr = nil
		mockfn.outParameters = nil
		mockfn.calls = nil
		mockfn.customArgs = nil
		mockfn = nil
	})

	Describe("newMockFunction", func() {
		var callCount int

		BeforeEach(func() {
			callCount = 0
			fn = func(str string, num int) (int, error) {
				callCount++
				return len(str) + num, nil
			}
		})

		It("returns error if passed a nil as the function pointer", func() {
			mockFn, err := newMockFunction(nil, nil)

			Expect(mockFn).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: expected the first argument to be a pointer to a function, but received a nil"))
		})

		It("returns error if a non-pointer value is passed as the function pointer", func() {
			mockFn, err := newMockFunction(42, nil)

			Expect(mockFn).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: expected the first argument to be a pointer to a function, but received a int"))
		})

		It("returns error if a non-function value is passed as the function pointer", func() {
			num := 42
			mockFn, err := newMockFunction(&num, nil)

			Expect(mockFn).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: expected the first argument to be a pointer to a function, but received a pointer to a int"))
		})

		It("returns error supplied out parameters are not of the same type", func() {
			mockFn, err := newMockFunction(&fn, []interface{}{"42", nil})

			Expect(mockFn).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: expected return values of type (int, error), but received (string, <nil>)"))
		})

		It("returns error if cloneValue returns an error", func() {
			_cloneValue = func(interface{}, interface{}) error {
				return errors.New("Ope")
			}
			defer func() {
				_cloneValue = cloneValue
			}()

			mockFn, err := newMockFunction(&fn, []interface{}{42, nil})

			Expect(mockFn).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: could not clone function pointer to new memory address: Ope"))
		})

		It("returns a mockFunction with a reference to the original function", func() {
			mockFn, err := newMockFunction(&fn, []interface{}{42, nil})

			Expect(err).To(BeNil())
			Expect(mockFn).ToNot(BeNil())
			Expect(mockFn.originalFunc).ToNot(BeNil())

			_, _ = mockFn.originalFunc.(func(str string, num int) (int, error))("", 0)

			Expect(callCount).To(Equal(1))
		})

		It("returns a mockFunction with properties initialized with zero values", func() {
			mockFn, err := newMockFunction(&fn, []interface{}{42, nil})

			Expect(err).To(BeNil())
			Expect(mockFn).ToNot(BeNil())
			Expect(mockFn.calls).To(BeNil())
			Expect(mockFn.customArgs).To(BeNil())
		})

		It("returns a mockFunction with outParameters as supplied", func() {
			mockFn, err := newMockFunction(&fn, []interface{}{42, nil})

			Expect(err).To(BeNil())
			Expect(mockFn).ToNot(BeNil())
			Expect(mockFn.outParameters).To(Equal([]interface{}{42, nil}))
		})
	})

	Describe("getReturnValues", func() {
		It("returns the mockFunction.OutParameters if no customArgs or onCalls exist", func() {
			args := []interface{}{"Hello", 42}

			result, maybeCustomArguments := mockfn.getReturnValues(args, reflect.TypeOf(fn))

			Expect(result).To(Equal([]interface{}{42, nil}))
			Expect(maybeCustomArguments).To(BeNil())
		})

		It("returns the mockFunction.OutParameters if the customArgs are nil", func() {
			args := []interface{}{"Hello", 42}
			mockfn.customArgs = append(mockfn.customArgs, nil, nil, nil)

			result, maybeCustomArguments := mockfn.getReturnValues(args, reflect.TypeOf(fn))

			Expect(result).To(Equal([]interface{}{42, nil}))
			Expect(maybeCustomArguments).To(BeNil())
		})

		It("returns the mockFunction.OutParameters if the argument do not equal any customArgs", func() {
			args := []interface{}{"Hello", 42}
			mockfn.customArgs = append(
				mockfn.customArgs,
				&customArguments{
					stub:        mockfn,
					argMatchers: []match.SupportedKindsMatcher{match.Exactly("apple"), match.IntGreaterThanOrEqualTo(0)},
					out:         []interface{}{0, errors.New("I am not an apple")},
				},
				&customArguments{
					stub:        mockfn,
					argMatchers: []match.SupportedKindsMatcher{match.Exactly("B"), match.Exactly(63)},
					out:         []interface{}{98, nil},
				})

			result, maybeCustomArguments := mockfn.getReturnValues(args, reflect.TypeOf(fn))

			Expect(result).To(Equal([]interface{}{42, nil}))
			Expect(maybeCustomArguments).To(BeNil())
		})

		It("returns the out parameters if the customArgs if the arguments match", func() {
			args := []interface{}{"Hello", 42}
			expected := &customArguments{
				stub:        mockfn,
				argMatchers: []match.SupportedKindsMatcher{match.StringSuffix("ello"), match.Exactly(42)},
				out:         []interface{}{22, errors.New("I am an error")},
			}
			mockfn.customArgs = append(
				mockfn.customArgs,
				&customArguments{
					stub:        mockfn,
					argMatchers: []match.SupportedKindsMatcher{match.StringPrefix("app"), match.Exactly(0)},
					out:         []interface{}{0, errors.New("I am not an apple")},
				},
				&customArguments{
					stub:        mockfn,
					argMatchers: []match.SupportedKindsMatcher{match.Exactly("B"), match.IntGreaterThan(60)},
					out:         []interface{}{98, nil},
				},
				expected,
			)

			result, maybeCustomArguments := mockfn.getReturnValues(args, reflect.TypeOf(fn))

			Expect(result).To(Equal([]interface{}{22, errors.New("I am an error")}))
			Expect(maybeCustomArguments).To(Equal(expected))
		})

		It("returns the the out parameters for the specific call index", func() {
			args := []interface{}{"apples", 42}
			mockfn.customArgs = append(
				mockfn.customArgs,
				&customArguments{
					stub:        mockfn,
					argMatchers: []match.SupportedKindsMatcher{match.Exactly("apple"), match.Exactly(0)},
					out:         []interface{}{0, errors.New("I am not an apple")},
					callCount:   2,
					onCalls: []*onCall{
						&onCall{
							stub:  mockfn,
							index: 2,
							out:   []interface{}{23, errors.New("I am the third not an apple")},
						},
					},
				},
			)
			mockfn.calls = []call{call{}, call{}, call{}}
			mockfn.onCalls = append(mockfn.onCalls, &onCall{
				stub:  mockfn,
				index: 3,
				out:   []interface{}{22, errors.New("I am the first error")},
			})

			result, maybeCustomArguments := mockfn.getReturnValues(args, reflect.TypeOf(fn))

			Expect(result).To(Equal([]interface{}{22, errors.New("I am the first error")}))
			Expect(maybeCustomArguments).To(BeNil())
		})

		It("returns the the out parameters for the specific call index on a custom argument", func() {
			args := []interface{}{"apple", 0}
			expected := &customArguments{
				stub:        mockfn,
				argMatchers: []match.SupportedKindsMatcher{match.Exactly("apple"), match.Exactly(0)},
				out:         []interface{}{0, errors.New("I am not an apple")},
				callCount:   2,
				onCalls: []*onCall{
					&onCall{
						stub:  mockfn,
						index: 2,
						out:   []interface{}{23, errors.New("I am the third not an apple")},
					},
				},
			}
			mockfn.customArgs = append(mockfn.customArgs, expected)
			mockfn.calls = []call{call{}, call{}, call{}}
			mockfn.onCalls = append(mockfn.onCalls, &onCall{
				stub:  mockfn,
				index: 0,
				out:   []interface{}{22, errors.New("I am the first error")},
			})

			result, maybeCustomArguments := mockfn.getReturnValues(args, reflect.TypeOf(fn))

			Expect(result).To(Equal([]interface{}{23, errors.New("I am the third not an apple")}))
			Expect(maybeCustomArguments).To(Equal(expected))
		})
	})

	Describe("toType", func() {
		It("returns the tuype of the mocked function", func() {
			fnValue := reflect.ValueOf(&fn).Elem()

			result := mockfn.toType()

			Expect(result).To(Equal(fnValue.Type()))
		})
	})

	Describe("implementation", func() {
		BeforeEach(func() {
			mockfn.outParameters = []interface{}{42, nil}
			mockfn.calls = []call{}
			mockfn.customArgs = []*customArguments{
				&customArguments{
					stub:        mockfn,
					argMatchers: []match.SupportedKindsMatcher{match.Exactly("custom"), match.Exactly(0)},
					out:         []interface{}{0, errors.New("Ope")},
				},
				&customArguments{
					stub:        mockfn,
					argMatchers: []match.SupportedKindsMatcher{match.StringPrefix("custom-"), match.Exactly(0)},
					out:         nil,
				},
			}
		})

		It("calls the exec function with the arguments as interfaces", func() {
			var argsProvided []interface{}
			args := []reflect.Value{reflect.ValueOf("Hello"), reflect.ValueOf(42)}
			expected := []interface{}{"Hello", 42}

			mockfn.execFunc = func(a []interface{}) {
				argsProvided = a
			}

			_ = mockfn.implementation(args)

			Expect(argsProvided).To(Equal(expected))
		})

		It("appends the call meta data for this call into the calls slice", func() {
			args := []reflect.Value{reflect.ValueOf("Hello"), reflect.ValueOf(42)}

			Expect(mockfn.calls).To(HaveLen(0))

			_ = mockfn.implementation(args)

			Expect(mockfn.calls).To(HaveLen(1))

			ca := mockfn.calls[0]

			Expect(ca.args).To(Equal([]interface{}{"Hello", 42}))
			Expect(ca.out).To(Equal([]interface{}{42, nil}))
		})

		It("returns the out parameters as reflection values", func() {
			args := []reflect.Value{reflect.ValueOf("Hello"), reflect.ValueOf(42)}

			outValues := mockfn.implementation(args)
			outInterfaces := mapToInterfaces(outValues)

			Expect(outInterfaces).To(Equal([]interface{}{42, nil}))
		})

		It("uses out parameters for custom arguments if applicable", func() {
			args := []reflect.Value{reflect.ValueOf("custom"), reflect.ValueOf(0)}

			outValues := mockfn.implementation(args)
			outInterfaces := mapToInterfaces(outValues)

			Expect(outInterfaces).To(Equal([]interface{}{0, errors.New("Ope")}))
		})

		It("doesn't use custom arguments if the return values have not been supplied", func() {
			args := []reflect.Value{reflect.ValueOf("custom-missing"), reflect.ValueOf(0)}

			outValues := mockfn.implementation(args)
			outInterfaces := mapToInterfaces(outValues)

			Expect(outInterfaces).To(Equal([]interface{}{42, nil}))
		})

		Context("variadic function", func() {
			BeforeEach(func() {
				fn := func(str string, opts ...string) (int, error) {
					return len(str) + len(opts), nil
				}

				mockfn.functionPtr = &fn
				Expect(cloneValue(&fn, &mockfn.originalFunc)).To(Succeed())
				mockfn.outParameters = []interface{}{42, nil}
				mockfn.calls = []call{}
			})

			It("appends the call meta data omitting the missing variadic arguments", func() {
				args := []reflect.Value{reflect.ValueOf("Hello")}

				Expect(mockfn.calls).To(HaveLen(0))

				_ = mockfn.implementation(args)

				Expect(mockfn.calls).To(HaveLen(1))

				ca := mockfn.calls[0]

				Expect(ca.args).To(Equal([]interface{}{"Hello"}))
			})

			It("appends the call meta data spreading the variadic arguments", func() {
				args := []reflect.Value{reflect.ValueOf("Hello"), reflect.ValueOf("A"), reflect.ValueOf("B")}

				Expect(mockfn.calls).To(HaveLen(0))

				_ = mockfn.implementation(args)

				Expect(mockfn.calls).To(HaveLen(1))

				ca := mockfn.calls[0]

				Expect(ca.args).To(Equal([]interface{}{"Hello", "A", "B"}))
			})
		})
	})

	Describe("Return", func() {
		It("returns error if the out parameters are not valid", func() {
			err := mockfn.Return(42, 42)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: expected return values of type (int, error), but received (int, int)"))
		})

		It("replaces the out parameters if they are valid", func() {
			err := mockfn.Return(22, errors.New("I am new"))

			Expect(err).To(BeNil())
			Expect(mockfn.outParameters).To(Equal([]interface{}{22, errors.New("I am new")}))
		})
	})

	Describe("WithArgs", func() {
		It("returns existing custom arguments if it matching arguments", func() {
			ca := &customArguments{
				stub:        mockfn,
				argMatchers: []match.SupportedKindsMatcher{match.Exactly("apple"), match.Exactly(0)},
				out:         []interface{}{0, errors.New("I am not an apple")},
			}
			mockfn.customArgs = append(mockfn.customArgs, ca)

			withArgs := mockfn.WithArgs("apple", 0)

			Expect(withArgs).ToNot(BeNil())
			Expect(withArgs).To(Equal(ca))
		})

		It("creates and returns new custom arguments if one does not exist", func() {
			Expect(mockfn.customArgs).To(HaveLen(0))

			withArgs := mockfn.WithArgs("apple", 0)

			Expect(withArgs).ToNot(BeNil())
			Expect(mockfn.customArgs).To(HaveLen(1))

			ca, ok := withArgs.(*customArguments)
			Expect(ok).To(BeTrue())

			Expect(ca.stub).To(Equal(mockfn))
			Expect(ca.argMatchers).To(Equal([]match.SupportedKindsMatcher{match.Exactly("apple"), match.Exactly(0)}))
			Expect(ca.out).To(BeNil())
		})
	})

	Describe("CallCount", func() {
		It("returns the number of times the stub was called", func() {
			Expect(mockfn.CallCount()).To(Equal(0))

			mockfn.calls = append(mockfn.calls, call{})

			Expect(mockfn.CallCount()).To(Equal(1))

			mockfn.calls = append(mockfn.calls, call{})

			Expect(mockfn.CallCount()).To(Equal(2))
		})
	})

	Describe("GetCalls", func() {
		It("returns all the call meta data made to the stub", func() {
			mockfn.calls = append(mockfn.calls,
				call{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				call{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				call{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			)

			result := mockfn.GetCalls()

			Expect(result).To(Equal([]Call{
				&call{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				&call{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				&call{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}))
		})
	})

	Describe("GetCall", func() {
		BeforeEach(func() {
			mockfn.calls = []call{
				call{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				call{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				call{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}
		})

		It("panics if the call index is less than 0", func() {
			defer func() {
				r := recover()
				if r == nil {
					Fail("expected a panic")
				}

				err, ok := r.(error)
				Expect(ok).To(BeTrue())
				Expect(err.Error()).To(Equal("mocka: attempted to get CallMetaData for call -1, when the function has only been called 3 times"))
			}()

			_ = mockfn.GetCall(-1)
			Fail("expected test to panic")
		})

		It("panics if the call index is greater than the number of calls", func() {
			defer func() {
				r := recover()
				if r == nil {
					Fail("expected a panic")
				}

				err, ok := r.(error)
				Expect(ok).To(BeTrue())
				Expect(err.Error()).To(Equal("mocka: attempted to get CallMetaData for call 5, when the function has only been called 3 times"))
			}()

			_ = mockfn.GetCall(5)
			Fail("expected test to panic")
		})

		It("returns the call meta data for the specified call index", func() {
			result := mockfn.GetCall(1)

			Expect(result).ToNot(BeNil())

			c := result.(*call)

			Expect(*c).To(Equal(call{
				args: []interface{}{"sam", 22},
				out:  []interface{}{42, nil},
			}))
		})
	})

	Describe("GetFirstCall", func() {
		It("panics if the stub has not been called once", func() {
			defer func() {
				r := recover()
				if r == nil {
					Fail("expected a panic")
				}

				err, ok := r.(error)
				Expect(ok).To(BeTrue())
				Expect(err.Error()).To(Equal("mocka: attempted to get CallMetaData for call 0, when the function has only been called 0 times"))
			}()

			_ = mockfn.GetFirstCall()
			Fail("expected test to panic")
		})

		It("returns the call meta data for the first call", func() {
			mockfn.calls = []call{
				call{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				call{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				call{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}

			result := mockfn.GetFirstCall()

			Expect(result).ToNot(BeNil())

			c := result.(*call)

			Expect(*c).To(Equal(call{
				args: []interface{}{"hello", 42},
				out:  []interface{}{22, nil},
			}))
		})
	})

	Describe("GetSecondCall", func() {
		It("panics if the stub has not been called twice", func() {
			defer func() {
				r := recover()
				if r == nil {
					Fail("expected a panic")
				}

				err, ok := r.(error)
				Expect(ok).To(BeTrue())
				Expect(err.Error()).To(Equal("mocka: attempted to get CallMetaData for call 1, when the function has only been called 0 times"))
			}()

			_ = mockfn.GetSecondCall()
			Fail("expected test to panic")
		})

		It("returns the call meta data for the second call", func() {
			mockfn.calls = []call{
				call{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				call{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				call{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}

			result := mockfn.GetSecondCall()

			Expect(result).ToNot(BeNil())

			c := result.(*call)

			Expect(*c).To(Equal(call{
				args: []interface{}{"sam", 22},
				out:  []interface{}{42, nil},
			}))
		})
	})

	Describe("GetThirdCall", func() {
		It("panics if the stub has not been called three times", func() {
			defer func() {
				r := recover()
				if r == nil {
					Fail("expected a panic")
				}

				err, ok := r.(error)
				Expect(ok).To(BeTrue())
				Expect(err.Error()).To(Equal("mocka: attempted to get CallMetaData for call 2, when the function has only been called 0 times"))
			}()

			_ = mockfn.GetThirdCall()
			Fail("expected test to panic")
		})

		It("returns the call meta data for the third call", func() {
			mockfn.calls = []call{
				call{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				call{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				call{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}

			result := mockfn.GetThirdCall()

			Expect(result).ToNot(BeNil())

			c := result.(*call)

			Expect(*c).To(Equal(call{
				args: []interface{}{"rob", 12},
				out:  []interface{}{0, errors.New("ope")},
			}))
		})
	})

	Describe("CalledOnce", func() {
		It("returns true is the stub has been called at least once", func() {
			mockfn.calls = []call{
				call{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				call{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				call{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}

			result := mockfn.CalledOnce()

			Expect(result).To(BeTrue())
		})

		It("returns false if the stub has not been called at least once", func() {
			result := mockfn.CalledOnce()

			Expect(result).To(BeFalse())
		})
	})

	Describe("CalledTwice", func() {
		It("returns true is the stub has been called at least twice", func() {
			mockfn.calls = []call{
				call{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				call{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				call{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}

			result := mockfn.CalledTwice()

			Expect(result).To(BeTrue())
		})

		It("returns false if the stub has not been called at least twice", func() {
			result := mockfn.CalledTwice()

			Expect(result).To(BeFalse())
		})
	})

	Describe("CalledThrice", func() {
		It("returns true is the stub has been called at least thrice", func() {
			mockfn.calls = []call{
				call{
					args: []interface{}{"hello", 42},
					out:  []interface{}{22, nil},
				},
				call{
					args: []interface{}{"sam", 22},
					out:  []interface{}{42, nil},
				},
				call{
					args: []interface{}{"rob", 12},
					out:  []interface{}{0, errors.New("ope")},
				},
			}

			result := mockfn.CalledThrice()

			Expect(result).To(BeTrue())
		})

		It("returns false if the stub has not been called at least thrice", func() {
			result := mockfn.CalledThrice()

			Expect(result).To(BeFalse())
		})
	})

	Describe("OnCall", func() {
		BeforeEach(func() {
			mockfn.onCalls = []*onCall{
				&onCall{stub: mockfn, index: 1},
				&onCall{stub: mockfn, index: 2},
				&onCall{stub: mockfn, index: 0},
			}
		})

		It("returns a pointer to an onCall struct", func() {
			result := mockfn.OnCall(5)

			o, ok := result.(*onCall)

			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockfn, index: 5}))
		})

		It("appends the new onCall struct to the onCalls slice", func() {
			_ = mockfn.OnCall(5)

			Expect(mockfn.onCalls).To(HaveLen(4))
		})

		It("returns an existing onCall object if one exists for that index", func() {
			result := mockfn.OnCall(2)

			o, ok := result.(*onCall)

			Expect(mockfn.onCalls).To(HaveLen(3))
			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockfn, index: 2}))
		})
	})

	Describe("OnFirstCall", func() {
		It("creates a new onCall with a 0 index", func() {
			mockfn.onCalls = []*onCall{
				&onCall{stub: mockfn, index: 1},
				&onCall{stub: mockfn, index: 2},
			}

			result := mockfn.OnFirstCall()

			o, ok := result.(*onCall)

			Expect(mockfn.onCalls).To(HaveLen(3))
			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockfn, index: 0}))
		})

		It("returns the existing first call", func() {
			mockfn.onCalls = []*onCall{
				&onCall{stub: mockfn, index: 1},
				&onCall{stub: mockfn, index: 2},
				&onCall{stub: mockfn, index: 0},
			}

			result := mockfn.OnFirstCall()

			o, ok := result.(*onCall)

			Expect(mockfn.onCalls).To(HaveLen(3))
			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockfn, index: 0}))
		})
	})

	Describe("OnSecondCall", func() {
		It("creates a new onCall with a 1 index", func() {
			mockfn.onCalls = []*onCall{
				&onCall{stub: mockfn, index: 0},
				&onCall{stub: mockfn, index: 2},
			}

			result := mockfn.OnSecondCall()

			o, ok := result.(*onCall)

			Expect(mockfn.onCalls).To(HaveLen(3))
			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockfn, index: 1}))
		})

		It("returns the existing second call", func() {
			mockfn.onCalls = []*onCall{
				&onCall{stub: mockfn, index: 1},
				&onCall{stub: mockfn, index: 2},
				&onCall{stub: mockfn, index: 0},
			}

			result := mockfn.OnSecondCall()

			o, ok := result.(*onCall)

			Expect(mockfn.onCalls).To(HaveLen(3))
			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockfn, index: 1}))
		})
	})

	Describe("OnThirdCall", func() {
		It("creates a new onCall with a 2 index", func() {
			mockfn.onCalls = []*onCall{
				&onCall{stub: mockfn, index: 0},
				&onCall{stub: mockfn, index: 1},
			}

			result := mockfn.OnThirdCall()

			o, ok := result.(*onCall)

			Expect(mockfn.onCalls).To(HaveLen(3))
			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockfn, index: 2}))
		})

		It("returns the existing third call", func() {
			mockfn.onCalls = []*onCall{
				&onCall{stub: mockfn, index: 1},
				&onCall{stub: mockfn, index: 2},
				&onCall{stub: mockfn, index: 0},
			}

			result := mockfn.OnThirdCall()

			o, ok := result.(*onCall)

			Expect(mockfn.onCalls).To(HaveLen(3))
			Expect(ok).To(BeTrue())
			Expect(*o).To(Equal(onCall{stub: mockfn, index: 2}))
		})
	})

	Describe("Restore", func() {
		It("overrides the function pointer with the the original pointer", func() {
			mockfn.originalFunc = func(str string, num int) (int, error) {
				return 42, errors.New("Ope")
			}

			castfn, ok := mockfn.functionPtr.(*func(str string, num int) (int, error))
			Expect(ok).To(BeTrue())

			n, err := (*castfn)("hello", 2)
			Expect(n).To(Equal(7))
			Expect(err).To(BeNil())

			mockfn.Restore()

			castfn, ok = mockfn.functionPtr.(*func(str string, num int) (int, error))
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

			mockfn.ExecOnCall(newFn)

			Expect(called).To(BeFalse())

			mockfn.execFunc([]interface{}{})

			Expect(called).To(BeTrue())
		})
	})

	Describe("getHighestPriority", func() {
		const numArguments = 2

		var (
			customArgs []*customArguments
			matcher1   *customArguments
			matcher2   *customArguments
		)

		BeforeEach(func() {
			matcher1 = &customArguments{
				stub:        mockfn,
				argMatchers: []match.SupportedKindsMatcher{match.Exactly("custom-"), match.Exactly(0)},
				out:         []interface{}{0, errors.New("Ope")},
			}
			matcher2 = &customArguments{
				stub:        mockfn,
				argMatchers: []match.SupportedKindsMatcher{match.StringPrefix("custom-"), match.Exactly(0)},
				out:         nil,
			}
			customArgs = []*customArguments{matcher1, matcher2}
		})

		It("returns nil when no custom arguments are provided", func() {
			actual := getHighestPriority(nil, numArguments)

			Expect(actual).To(BeNil())
		})

		It("return the only custom argument when provided one custom argument", func() {
			actual := getHighestPriority([]*customArguments{matcher1}, numArguments)

			Expect(actual).To(Equal(matcher1))
		})

		It("returns the matcher with the highest priority", func() {
			actual := getHighestPriority(customArgs, numArguments)

			Expect(actual).To(Equal(matcher1))
		})

		It("returns the first matcher if multiple matchers have the same priority", func() {
			customArgs = []*customArguments{
				matcher2,
				&customArguments{
					stub:        mockfn,
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
			customArgs []*customArguments
			matcher1   *customArguments
			matcher2   *customArguments
		)

		BeforeEach(func() {
			matcher1 = &customArguments{
				stub:        mockfn,
				argMatchers: []match.SupportedKindsMatcher{match.Exactly("custom-"), match.Exactly(0)},
				out:         []interface{}{0, errors.New("Ope")},
			}
			matcher2 = &customArguments{
				stub:        mockfn,
				argMatchers: []match.SupportedKindsMatcher{match.StringPrefix("custom-"), match.Exactly(0)},
				out:         nil,
			}
			customArgs = []*customArguments{matcher1, matcher2}
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
