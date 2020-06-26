package mocka

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mocka", func() {
	Describe("Function", func() {
		var (
			callCount int
			fn        func(str string, num int) (int, error)
		)

		BeforeEach(func() {
			callCount = 0
			fn = func(str string, num int) (int, error) {
				callCount++
				return len(str) + num, nil
			}
		})

		It("returns error if passed a nil as the function pointer", func() {
			stub, err := Function(nil)

			Expect(stub).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: expected the first argument to be a pointer to a function, but received a nil"))
		})

		It("returns error if a non-pointer value is passed as the function pointer", func() {
			stub, err := Function(42)

			Expect(stub).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: expected the first argument to be a pointer to a function, but received a int"))
		})

		It("returns error if a non-function value is passed as the function pointer", func() {
			num := 42
			stub, err := Function(&num)

			Expect(stub).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: expected the first argument to be a pointer to a function, but received a pointer to a int"))
		})

		It("returns error supplied out parameters are not of the same type", func() {
			stub, err := Function(&fn, "42", nil)

			Expect(stub).To(BeNil())
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

			stub, err := Function(&fn, 42, nil)

			Expect(stub).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("mocka: could not clone function pointer to new memory address: Ope"))
		})

		It("returns a stub with a reference to the original function", func() {
			stub, err := Function(&fn, 42, nil)

			Expect(err).To(BeNil())
			Expect(stub).ToNot(BeNil())

			Expect(stub.originalFunc).ToNot(BeNil())

			_, _ = stub.originalFunc.(func(str string, num int) (int, error))("", 0)

			Expect(callCount).To(Equal(1))
		})

		It("returns a stub with properties initialized with zero values", func() {
			stub, err := Function(&fn, 42, nil)

			Expect(err).To(BeNil())
			Expect(stub).ToNot(BeNil())
			Expect(stub.calls).To(BeNil())
			Expect(stub.customArgs).To(BeNil())
		})

		It("returns a stub with outParameters as supplied", func() {
			stub, err := Function(&fn, 42, nil)

			Expect(err).To(BeNil())
			Expect(stub).ToNot(BeNil())
			Expect(stub.outParameters).To(Equal([]interface{}{42, nil}))
		})
	})

	Describe("CreateSandbox", func() {
		It("returns a sandbox with stub assigned as nil", func() {
			s := CreateSandbox()

			Expect(s).ToNot(BeNil())
			Expect(s.stubs).To(BeNil())
		})
	})
})
