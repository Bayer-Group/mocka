package mocka

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Namer interface {
	Name() string
}

type Thing struct {
	name string
}

func (thing *Thing) Name() string {
	return thing.name
}

var _ = Describe("utils", func() {
	Describe("mapToInterfaces", func() {
		var (
			thing  = Thing{"The Thing"}
			values []reflect.Value
		)

		BeforeEach(func() {
			values = []reflect.Value{
				reflect.ValueOf("Hello"),
				reflect.ValueOf(42),
				reflect.ValueOf(thing),
				reflect.ValueOf(true),
				reflect.ValueOf(nil),
			}
		})

		It("returns the same values as interfaces", func() {
			result := mapToInterfaces(values)
			Expect(result).To(Equal([]interface{}{"Hello", 42, thing, true, nil}))
		})
	})

	Describe("mapToReflectValue", func() {
		var (
			thing      = Thing{"The Thing"}
			interfaces []interface{}
		)

		BeforeEach(func() {
			interfaces = []interface{}{"Hello", 42, thing, true, nil}
		})

		It("returns the same values as interfaces", func() {
			result := mapToReflectValue(interfaces)
			Expect(mapToInterfaces(result)).To(Equal(interfaces))
		})
	})

	Describe("cloneValue", func() {
		var aThing Thing

		BeforeEach(func() {
			aThing = Thing{name: "Jon"}
		})

		It("should create a deep clone of a struct pointer", func() {
			var bThing Thing

			err := cloneValue(&aThing, &bThing)

			Expect(err).To(BeNil())
			Expect(bThing).To(Equal(aThing))
			Expect(bThing.name).To(Equal(aThing.name))

			aThing.name = "John"

			Expect(aThing.name).To(Equal("John"))
			Expect(bThing.name).To(Equal("Jon"))
		})

		It("throws error if source is not a pointer", func() {
			var bThing Thing

			err := cloneValue(aThing, &bThing)
			Expect(err).To(Equal(fmt.Errorf("mocka: expected source value for clone to be a pointer, but it was a struct")))
		})

		It("throws error if destination is not a pointer", func() {
			var bThing Thing

			err := cloneValue(&aThing, bThing)
			Expect(err).To(Equal(fmt.Errorf("mocka: expected destination value for clone to be a pointer, but it was a struct")))
		})
	})

	Describe("validateOutParameters", func() {
		var fn = func(str string, num int) (int, error) {
			return len(str) + num, nil
		}

		It("returns false if the type is nil", func() {
			result := validateOutParameters(nil, []interface{}{0, nil})
			Expect(result).To(BeFalse())
		})

		It("returns false if the type is not a function", func() {
			result := validateOutParameters(reflect.TypeOf("I am not a funciton"), []interface{}{0, nil})
			Expect(result).To(BeFalse())
		})

		It("returns false of the number of out parameters /= number of function out parameters", func() {
			result := validateOutParameters(reflect.TypeOf(fn), []interface{}{0, nil, 42})
			Expect(result).To(BeFalse())
		})

		It("returns false if one of the out parameters type does not match", func() {
			result := validateOutParameters(reflect.TypeOf(fn), []interface{}{"0", errors.New("I am an error")})
			Expect(result).To(BeFalse())
		})

		It("returns true if all out parameter types match", func() {
			result := validateOutParameters(reflect.TypeOf(fn), []interface{}{42, errors.New("I am an error")})
			Expect(result).To(BeTrue())
		})
	})

	Describe("areTypeAndValueEquivalent", func() {
		It("returns false if the type is nil", func() {
			Expect(areTypeAndValueEquivalent(nil, "")).To(BeFalse())
		})

		It("returns true if nil for type is a nilable kind", func() {
			var namer Namer
			namer = &Thing{"The Thing"}
			num := 0
			nilables := map[reflect.Type]interface{}{
				reflect.TypeOf(errors.New("")):          nil,
				reflect.TypeOf(func() {}):               (func())(nil),
				reflect.TypeOf(namer):                   Namer(nil),
				reflect.TypeOf(make(chan int)):          (chan int)(nil),
				reflect.TypeOf(make(map[string]string)): (map[string]string)(nil),
				reflect.TypeOf(&num):                    (*int)(nil),
				reflect.TypeOf([]int{}):                 ([]int)(nil),
				reflect.TypeOf(fmt.Errorf("")):          (error)(nil),
			}

			for valueType, value := range nilables {
				Expect(areTypeAndValueEquivalent(valueType, value)).To(BeTrue())
			}
		})

		It("returns true if initialized value for type is a nilable kind", func() {
			var namer Namer
			namer = &Thing{"The Thing"}
			num := 0
			nilables := map[reflect.Type]interface{}{
				reflect.TypeOf(errors.New("")):          errors.New("I am an error"),
				reflect.TypeOf(namer):                   namer,
				reflect.TypeOf(func() {}):               func() {},
				reflect.TypeOf(make(chan int)):          make(chan int, 10),
				reflect.TypeOf(make(map[string]string)): map[string]string{"key": "value"},
				reflect.TypeOf(&num):                    &num,
				reflect.TypeOf([]int{}):                 []int{1, 2, 3},
			}

			for valueType, value := range nilables {
				Expect(areTypeAndValueEquivalent(valueType, value)).To(BeTrue())
			}
		})

		It("return true if non-nil kind's match", func() {
			nilables := map[reflect.Type]interface{}{
				reflect.TypeOf((string)("")):   "adf",
				reflect.TypeOf((int)(1)):       150,
				reflect.TypeOf((float64)(1.5)): 30.2,
			}

			for valueType, value := range nilables {
				Expect(areTypeAndValueEquivalent(valueType, value)).To(BeTrue())
			}
		})

		It("return false if non-nil kind's do not match", func() {
			nilables := map[reflect.Type]interface{}{
				reflect.TypeOf((string)("")):   150,
				reflect.TypeOf((int)(1)):       "asdf",
				reflect.TypeOf((float64)(1.5)): 30,
			}

			for valueType, value := range nilables {
				Expect(areTypeAndValueEquivalent(valueType, value)).To(BeFalse())
			}
		})
	})

	Describe("mapToTypeName", func() {
		type thisIsAStruct struct{}

		It("returns the type name of each element in a slice", func() {
			str := "pointer"
			thing := &Thing{}
			namer := Namer(thing)
			input := []interface{}{namer, 10, nil, "Hello", 10.0, errors.New("Ope"), thisIsAStruct{}, &str}

			result := mapToTypeName(input)

			Expect(result).To(Equal([]string{"*Thing", "int", "<nil>", "string", "float64", "*errorString", "thisIsAStruct", "*string"}))
		})

		It("returns an empty slice if passed a nil", func() {
			Expect(mapToTypeName(nil)).To(Equal([]string{}))
		})
	})
})
