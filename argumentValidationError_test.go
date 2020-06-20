package mocka

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("argumentValidationError", func() {
	var fnType reflect.Type

	BeforeEach(func() {
		var fn = func(str string, num int) (int, error) {
			return len(str) + num, nil
		}
		fnType = reflect.ValueOf(&fn).Elem().Type()
	})

	Describe("String", func() {
		It("returns a less descriptive error if fnType is nil", func() {
			arguments := []interface{}{0, ""}
			err := &argumentValidationError{nil, arguments}

			result := err.String()

			Expect(result).To(Equal("mocka: expected arguments of (int, string) to match function arguments"))
		})

		It("returns the error string with the type names of the expected and actual arguments", func() {
			arguments := []interface{}{0, ""}
			err := &argumentValidationError{fnType, arguments}

			result := err.String()

			Expect(result).To(Equal("mocka: expected arguments of type (string, int), but received (int, string)"))
		})

		It("returns the error string using ... to denote variadic arguments", func() {
			arguments := []interface{}{0, 0}
			var fn = func(str string, opts ...string) int {
				return len(str) + len(opts)
			}
			fnType = reflect.ValueOf(&fn).Elem().Type()

			err := &argumentValidationError{fnType, arguments}
			result := err.String()

			Expect(result).To(Equal("mocka: expected arguments of type (string, ...string), but received (int, int)"))
		})
	})

	Describe("Error", func() {
		It("returns a less descriptive error if fnType is nil", func() {
			arguments := []interface{}{0, ""}
			err := &argumentValidationError{nil, arguments}

			result := err.String()

			Expect(result).To(Equal("mocka: expected arguments of (int, string) to match function arguments"))
		})

		It("returns the error string with the type names of the expected and actual arguments", func() {
			arguments := []interface{}{0, ""}
			err := &argumentValidationError{fnType, arguments}

			result := err.Error()

			Expect(result).To(Equal("mocka: expected arguments of type (string, int), but received (int, string)"))
		})
	})
})
