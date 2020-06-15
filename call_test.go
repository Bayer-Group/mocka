package mocka

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("call", func() {
	Describe("newCall", func() {
		It("returns a call with the provided arguments and return values", func() {
			var fn func(string, string) int
			fnType := reflect.TypeOf(fn)

			testCall := newCall(fnType, []interface{}{"a", "b"}, []interface{}{1})

			Expect(testCall).To(Equal(call{
				args: []interface{}{"a", "b"},
				out:  []interface{}{1},
			}))
		})

		Context("variadic function", func() {
			var (
				fnType reflect.Type
				out    = []interface{}{1}
			)

			BeforeEach(func() {
				var fn func(string, ...string) int
				fnType = reflect.TypeOf(fn)
			})

			It("returns only non-variadic arguments if no variadic arguments are provided", func() {
				testCall := newCall(fnType, []interface{}{"a"}, out)

				Expect(testCall).To(Equal(call{
					args: []interface{}{"a"},
					out:  out,
				}))
			})

			It("flattens variadic arguments into call.args", func() {
				testCall := newCall(fnType, []interface{}{"a", []string{"b", "c", "d"}}, out)

				Expect(testCall).To(Equal(call{
					args: []interface{}{"a", "b", "c", "d"},
					out:  out,
				}))
			})
		})
	})

	Describe("Arguments", func() {
		It("returns the arguments of the call", func() {
			testCall := &call{
				args: []interface{}{42, "hello"},
				out:  []interface{}{40, nil},
			}

			result := testCall.Arguments()

			Expect(result).To(Equal([]interface{}{42, "hello"}))
		})
	})

	Describe("ReturnValues", func() {
		It("returns the return values of the call", func() {
			testCall := &call{
				args: []interface{}{42, "hello"},
				out:  []interface{}{40, nil},
			}

			result := testCall.ReturnValues()

			Expect(result).To(Equal([]interface{}{40, nil}))
		})
	})
})
