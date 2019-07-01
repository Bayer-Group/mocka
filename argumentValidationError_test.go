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

	Describe("Error", func() {
		It("returns the error string with the type names of the expected and actual arguments", func() {
			arguments := []interface{}{0, ""}
			err := &argumentValidationError{fnType, arguments}

			result := err.Error()

			Expect(result).To(Equal("mocka: expected arguments of type [string int], but recieved [int string]"))
		})
	})
})
