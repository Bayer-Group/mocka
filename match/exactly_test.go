package match

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("exactly", func() {
	var (
		mockChan1    chan int
		mockChan2    chan string
		mockFn1      func()
		mockFn2      func(_ string)
		mockPointer1 *exactly
		mockPointer2 *anything
	)

	BeforeEach(func() {
		mockChan1 = make(chan int, 1)
		mockChan2 = make(chan string, 1)
		mockFn1 = func() {}
		mockFn2 = func(_ string) {}
		mockPointer1 = &exactly{123}
		mockPointer2 = &anything{}
	})

	Describe("Exactly", func() {
		It("returns an exactly struct", func() {
			actual := Exactly(nil)

			Expect(actual).To(BeAssignableToTypeOf(new(exactly)))
		})
	})

	Describe("SupportedKinds", func() {
		It("returns all support kinds in go", func() {
			actual := Exactly(123).SupportedKinds()

			Expect(actual).To(Equal(
				map[reflect.Kind]struct{}{
					reflect.Bool:          struct{}{},
					reflect.Int:           struct{}{},
					reflect.Int8:          struct{}{},
					reflect.Int16:         struct{}{},
					reflect.Int32:         struct{}{},
					reflect.Int64:         struct{}{},
					reflect.Uint:          struct{}{},
					reflect.Uint8:         struct{}{},
					reflect.Uint16:        struct{}{},
					reflect.Uint32:        struct{}{},
					reflect.Uint64:        struct{}{},
					reflect.Uintptr:       struct{}{},
					reflect.Float32:       struct{}{},
					reflect.Float64:       struct{}{},
					reflect.Complex64:     struct{}{},
					reflect.Complex128:    struct{}{},
					reflect.Array:         struct{}{},
					reflect.Chan:          struct{}{},
					reflect.Func:          struct{}{},
					reflect.Interface:     struct{}{},
					reflect.Map:           struct{}{},
					reflect.Ptr:           struct{}{},
					reflect.Slice:         struct{}{},
					reflect.String:        struct{}{},
					reflect.Struct:        struct{}{},
					reflect.UnsafePointer: struct{}{},
				}))
		})
	})

	DescribeTable("Match",
		func(first interface{}, second interface{}, expected bool) {
			if expected {
				Expect(Exactly(first).Match(second)).To(BeTrue())
			} else {
				Expect(Exactly(first).Match(second)).To(BeFalse())
			}
		},
		Entry("returns true when both are nils", nil, nil, true),
		Entry("returns true when numbers are equal", 123, 123, true),
		Entry("returns true when bools are equal", true, true, true),
		Entry("returns true when structs are equal", exactly{123}, exactly{123}, true),
		Entry("returns true when channels are equal", mockChan1, mockChan1, true),
		Entry("returns true when functions are equal", mockFn1, mockFn1, true),
		Entry("returns true when pointers are equal", mockPointer1, mockPointer1, true),
		Entry("returns true when slices are equal", []string{"screams"}, []string{"screams"}, true),
		Entry("returns true when arrays are equal", [1]int{1}, [1]int{1}, true),
		Entry("returns true when maps are equal", map[string]struct{}{"a": struct{}{}}, map[string]struct{}{"a": struct{}{}}, true),
		Entry("returns false when numbers are not equal", 123, 563, false),
		Entry("returns false when bools are not equal", true, false, false),
		Entry("returns false when structs are not equal", exactly{"a"}, exactly{123}, false),
		Entry("returns false when channels are not equal", mockChan1, mockChan2, false),
		Entry("returns false when functions are not equal", mockFn1, mockFn2, false),
		Entry("returns false when pointers are not equal", mockPointer1, mockPointer2, false),
		Entry("returns false when slices are not equal", []string{"scremas"}, []string{"apple"}, false),
		Entry("returns false when arrays are not equal", [1]int{1}, [1]int{3}, false),
		Entry("returns false when maps are not equal", map[string]struct{}{"a": struct{}{}}, map[string]struct{}{"b": struct{}{}}, false),
	)
})
