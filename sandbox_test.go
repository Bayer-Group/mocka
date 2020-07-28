package mocka

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("sandbox", func() {
	var (
		callCounts       map[string]int
		fn1              func(string, int) (int, error)
		fn2              func(string) int
		fn3              func(interface{}) error
		testSandbox      *Sandbox
		failTestReporter *mockTestReporter
	)

	BeforeEach(func() {
		failTestReporter = &mockTestReporter{}
		testSandbox = &Sandbox{testReporter: GinkgoT()}
		callCounts = map[string]int{"fn1": 0, "fn2": 0, "fn3": 0}
		fn1 = func(str string, num int) (int, error) {
			callCounts["fn1"]++
			return len(str) + num, nil
		}
		fn2 = func(str string) int {
			callCounts["fn2"]++
			return len(str)
		}
		fn3 = func(i interface{}) error {
			callCounts["fn3"]++
			if i == nil {
				return errors.New("data is nil")
			}

			return nil
		}
	})

	AfterEach(func() {
		// clear out slice, to prevent memory leaks
		testSandbox.stubs = nil
	})

	Describe("Function", func() {
		It("reports an error if passed a nil as the function pointer", func() {
			testSandbox.testReporter = failTestReporter
			stub := testSandbox.Function(nil)

			Expect(stub).To(BeNil())
			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected the second argument to be a pointer to a function, but received a nil",
			}))
		})

		It("reports an error if a non-pointer value is passed as the function pointer", func() {
			testSandbox.testReporter = failTestReporter
			stub := testSandbox.Function(42)

			Expect(stub).To(BeNil())
			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected the second argument to be a pointer to a function, but received a int",
			}))
		})

		It("reports an error if a non-function value is passed as the function pointer", func() {
			testSandbox.testReporter = failTestReporter
			num := 42
			stub := testSandbox.Function(&num)

			Expect(stub).To(BeNil())
			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected the second argument to be a pointer to a function, but received a pointer to a int",
			}))
		})

		It("reports an error supplied out parameters are not of the same type", func() {
			testSandbox.testReporter = failTestReporter
			stub := testSandbox.Function(&fn1, "42", nil)

			Expect(stub).To(BeNil())
			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: expected return values of type (int, error), but received (string, <nil>)",
			}))
		})

		It("reports an error if cloneValue returns an error", func() {
			testSandbox.testReporter = failTestReporter
			_cloneValue = func(interface{}, interface{}) error {
				return errors.New("Ope")
			}
			defer func() {
				_cloneValue = cloneValue
			}()

			stub := testSandbox.Function(&fn1, 42, nil)

			Expect(stub).To(BeNil())
			Expect(failTestReporter.messages).To(Equal([]string{
				"mocka: could not clone function pointer to new memory address: Ope",
			}))
		})

		It("returns a stub with a reference to the original function", func() {
			stub := testSandbox.Function(&fn1, 42, nil)

			Expect(stub).ToNot(BeNil())

			Expect(stub.originalFunc).ToNot(BeNil())

			_, _ = stub.originalFunc.(func(str string, num int) (int, error))("", 0)

			Expect(callCounts["fn1"]).To(Equal(1))
		})

		It("returns a stub with properties initialized with zero values", func() {
			stub := testSandbox.Function(&fn1, 42, nil)

			Expect(stub).ToNot(BeNil())
			Expect(stub.calls).To(BeNil())
			Expect(stub.customArgs).To(BeNil())
		})

		It("returns a stub with outParameters as supplied", func() {
			stub := testSandbox.Function(&fn1, 42, nil)

			Expect(stub).ToNot(BeNil())
			Expect(stub.outParameters).To(Equal([]interface{}{42, nil}))
		})

		It("appends the stub into the sandbox if no error is returned", func() {
			Expect(testSandbox.stubs).To(HaveLen(0))

			stub := testSandbox.Function(&fn1, 42, nil)

			Expect(stub).ToNot(BeNil())
			Expect(testSandbox.stubs).To(HaveLen(1))
		})
	})

	Describe("Restore", func() {
		BeforeEach(func() {
			_ = testSandbox.Function(&fn1, 42, nil)
			_ = testSandbox.Function(&fn2, 42)
			_ = testSandbox.Function(&fn3, nil)
		})

		It("restores each function back to it's original functionality", func() {
			_, _ = fn1("", 0)
			_ = fn2("")
			_ = fn3(nil)

			Expect(callCounts["fn1"]).To(Equal(0))
			Expect(callCounts["fn2"]).To(Equal(0))
			Expect(callCounts["fn3"]).To(Equal(0))

			testSandbox.Restore()

			_, _ = fn1("", 0)
			_ = fn2("")
			_ = fn3(nil)

			Expect(callCounts["fn1"]).To(Equal(1))
			Expect(callCounts["fn2"]).To(Equal(1))
			Expect(callCounts["fn3"]).To(Equal(1))
		})

		It("removes references to the created Stubs", func() {
			Expect(testSandbox.stubs).To(HaveLen(3))

			testSandbox.Restore()

			Expect(testSandbox.stubs).To(HaveLen(0))
		})
	})
})
