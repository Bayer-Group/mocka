package mocka

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("call", func() {
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
