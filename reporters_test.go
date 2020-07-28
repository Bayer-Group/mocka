package mocka

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("reporters", func() {
	var (
		reporter     *mockTestReporter
		functionType reflect.Type
	)

	BeforeEach(func() {
		reporter = &mockTestReporter{}

		fn := func(str string, num int) (int, error) {
			return len(str) + num, nil
		}
		functionType = reflect.ValueOf(&fn).Elem().Type()
	})

	Describe("reportInvalidArguments", func() {
		It("reports the error string with the type names of the expected and actual arguments", func() {
			arguments := []interface{}{0, ""}

			reportInvalidArguments(reporter, functionType, arguments)

			Expect(reporter.messages).To(HaveLen(1))
			Expect(reporter.messages).To(ContainElement("mocka: expected arguments of type (string, int), but received (int, string)"))
		})

		It("reports the error string using ... to denote variadic arguments", func() {
			arguments := []interface{}{0, 0}
			var fn = func(str string, opts ...string) int {
				return len(str) + len(opts)
			}
			functionType = reflect.ValueOf(&fn).Elem().Type()

			reportInvalidArguments(reporter, functionType, arguments)

			Expect(reporter.messages).To(HaveLen(1))
			Expect(reporter.messages).To(ContainElement("mocka: expected arguments of type (string, ...string), but received (int, int)"))
		})
	})

	Describe("reportInvalidOutParameters", func() {
		It("reports a less descriptive error if fnType is nil", func() {
			outParameters := []interface{}{0, ""}

			reportInvalidOutParameters(reporter, nil, outParameters)

			Expect(reporter.messages).To(HaveLen(1))
			Expect(reporter.messages).To(ContainElement("mocka: expected return values of (int, string) to match function return values"))
		})

		It("reports the error string with the type names of the expected and actual return values", func() {
			outParameters := []interface{}{0, ""}

			reportInvalidOutParameters(reporter, functionType, outParameters)

			Expect(reporter.messages).To(HaveLen(1))
			Expect(reporter.messages).To(ContainElement("mocka: expected return values of type (int, error), but received (int, string)"))
		})
	})
})
