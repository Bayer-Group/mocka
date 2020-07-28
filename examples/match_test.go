package examples

import (
	"fmt"
	"reflect"

	"github.com/MonsantoCo/mocka/v2"
	"github.com/MonsantoCo/mocka/v2/match"
)

func ExampleAnything() {
	var fn = func(str []string, n int) int {
		return len(str) + n
	}

	stub := mocka.Function(t, &fn, 20)
	defer stub.Restore()

	stub.WithArgs(match.Anything(), 2).Return(10)
	stub.WithArgs([]string{"123", "456"}, 2).Return(5)

	fmt.Println(fn([]string{"hello"}, 5))
	fmt.Println(fn([]string{"mocka"}, 2))
	fmt.Println(fn([]string{"123", "456"}, 2))
	// Output: 20
	// 10
	// 5
}

func ExampleAnythingButNil() {
	var fn = func(str []string, n int) int {
		return len(str) + n
	}

	stub := mocka.Function(t, &fn, 20)
	defer stub.Restore()

	stub.WithArgs(match.AnythingButNil(), 2).Return(10)
	stub.WithArgs([]string{"123", "456"}, 2).Return(5)

	fmt.Println(fn([]string{"hello"}, 5))
	fmt.Println(fn([]string{"mocka"}, 2))
	fmt.Println(fn(nil, 2))
	fmt.Println(fn([]string{"123", "456"}, 2))
	// Output: 20
	// 10
	// 20
	// 5
}

func ExampleConvertibleTo() {
	var fn = func(x int, y int) int {
		return x + y
	}

	stub := mocka.Function(t, &fn, 20)
	defer stub.Restore()

	stub.WithArgs(match.ConvertibleTo(new(int64)), 2).Return(10)
	stub.WithArgs(10, 2).Return(5)

	fmt.Println(fn(10, 2))
	fmt.Println(fn(8, 2))
	fmt.Println(fn(8, 3))
	// Output: 5
	// 10
	// 20
}

func ExampleElementsContaining() {
	var fn = func(s []string) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.ElementsContaining("A", "B")).Return(30)
	stub.WithArgs(match.ElementsContaining("A")).Return(20)

	fmt.Println(fn([]string{"C"}))
	fmt.Println(fn([]string{"A", "C"}))
	fmt.Println(fn([]string{"A", "B", "C"}))
	// Output: 10
	// 20
	// 30
}

func ExampleEmpty() {
	var fn = func(s []string, m map[string]struct{}, str string) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.Empty(), match.Empty(), match.Empty()).Return(20)
	stub.WithArgs(match.Empty(), match.Empty(), "hello").Return(30)

	fmt.Println(fn(nil, nil, ""))
	fmt.Println(fn(nil, nil, "hello"))
	fmt.Println(fn([]string{""}, map[string]struct{}{"": {}}, "screams"))
	// Output: 20
	// 30
	// 10
}

func ExampleExactly() {
	var fn = func(x int) int {
		return x
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.Exactly(2)).Return(20)
	stub.WithArgs(match.Exactly(10)).Return(30)

	fmt.Println(fn(50))
	fmt.Println(fn(2))
	fmt.Println(fn(10))
	// Output: 10
	// 20
	// 30
}

func ExampleFloatGreaterThan() {
	var fn = func(x float64) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.FloatGreaterThan(2)).Return(20)

	fmt.Println(fn(1))
	fmt.Println(fn(3))
	// Output: 10
	// 20
}

func ExampleFloatGreaterThanOrEqualTo() {
	var fn = func(x float64) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.FloatGreaterThanOrEqualTo(2)).Return(20)

	fmt.Println(fn(1))
	fmt.Println(fn(2))
	fmt.Println(fn(3))
	// Output: 10
	// 20
	// 20
}

func ExampleFloatLessThan() {
	var fn = func(x float64) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.FloatLessThan(2)).Return(20)

	fmt.Println(fn(1))
	fmt.Println(fn(3))
	// Output: 20
	// 10
}

func ExampleFloatLessThanOrEqualTo() {
	var fn = func(x float64) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.FloatLessThanOrEqualTo(2)).Return(20)

	fmt.Println(fn(1))
	fmt.Println(fn(2))
	fmt.Println(fn(3))
	// Output: 20
	// 20
	// 10
}

func ExampleIntGreaterThan() {
	var fn = func(x int64) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.IntGreaterThan(2)).Return(20)

	fmt.Println(fn(1))
	fmt.Println(fn(3))
	// Output: 10
	// 20
}

func ExampleIntGreaterThanOrEqualTo() {
	var fn = func(x int64) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.IntGreaterThanOrEqualTo(2)).Return(20)

	fmt.Println(fn(1))
	fmt.Println(fn(2))
	fmt.Println(fn(3))
	// Output: 10
	// 20
	// 20
}

func ExampleIntLessThan() {
	var fn = func(x int64) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.IntLessThan(2)).Return(20)

	fmt.Println(fn(1))
	fmt.Println(fn(3))
	// Output: 20
	// 10
}

func ExampleIntLessThanOrEqualTo() {
	var fn = func(x int64) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.IntLessThanOrEqualTo(2)).Return(20)

	fmt.Println(fn(1))
	fmt.Println(fn(2))
	fmt.Println(fn(3))
	// Output: 20
	// 20
	// 10
}

func ExampleUintGreaterThan() {
	var fn = func(x uint64) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.UintGreaterThan(2)).Return(20)

	fmt.Println(fn(1))
	fmt.Println(fn(3))
	// Output: 10
	// 20
}

func ExampleUintGreaterThanOrEqualTo() {
	var fn = func(x uint64) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.UintGreaterThanOrEqualTo(2)).Return(20)

	fmt.Println(fn(1))
	fmt.Println(fn(2))
	fmt.Println(fn(3))
	// Output: 10
	// 20
	// 20
}

func ExampleUintLessThan() {
	var fn = func(x uint64) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.UintLessThan(2)).Return(20)

	fmt.Println(fn(1))
	fmt.Println(fn(3))
	// Output: 20
	// 10
}

func ExampleUintLessThanOrEqualTo() {
	var fn = func(x uint64) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.UintLessThanOrEqualTo(2)).Return(20)

	fmt.Println(fn(1))
	fmt.Println(fn(2))
	fmt.Println(fn(3))
	// Output: 20
	// 20
	// 10
}

func ExampleKeysContaining() {
	var fn = func(m map[string]struct{}) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.KeysContaining("A", "B")).Return(30)
	stub.WithArgs(match.KeysContaining("A")).Return(20)

	fmt.Println(fn(map[string]struct{}{"C": {}}))
	fmt.Println(fn(map[string]struct{}{"A": {}}))
	fmt.Println(fn(map[string]struct{}{"A": {}, "B": {}}))
	// Output: 10
	// 20
	// 30
}

func ExampleLengthOf() {
	var fn = func(s []int, m map[string]struct{}, str string) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.LengthOf(1), match.LengthOf(1), match.LengthOf(1)).Return(20)
	stub.WithArgs(match.LengthOf(1), match.LengthOf(1), match.LengthOf(2)).Return(30)

	fmt.Println(fn([]int{1, 2}, map[string]struct{}{"-": {}}, "-"))
	fmt.Println(fn([]int{1}, map[string]struct{}{"-": {}}, "-"))
	fmt.Println(fn([]int{1}, map[string]struct{}{"-": {}}, "--"))
	// Output: 10
	// 20
	// 30
}

func ExampleNil() {
	var fn = func(m map[string]struct{}) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.Nil()).Return(20)

	fmt.Println(fn(map[string]struct{}{"-": {}}))
	fmt.Println(fn(nil))
	// Output: 10
	// 20
}

func ExampleStringContaining() {
	var fn = func(s string) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.StringContaining("cream")).Return(20)

	fmt.Println(fn("apples"))
	fmt.Println(fn("screams"))
	// Output: 10
	// 20
}

func ExampleStringPrefix() {
	var fn = func(s string) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.StringPrefix("scr")).Return(20)

	fmt.Println(fn("apples"))
	fmt.Println(fn("screams"))
	// Output: 10
	// 20
}

func ExampleStringSuffix() {
	var fn = func(s string) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.StringSuffix("ms")).Return(20)

	fmt.Println(fn("apples"))
	fmt.Println(fn("screams"))
	// Output: 10
	// 20
}

func ExampleTypeOf() {
	var fn = func(s interface{}) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.TypeOf("int")).Return(20)
	stub.WithArgs(match.TypeOf("string")).Return(30)

	fmt.Println(fn(nil))
	fmt.Println(fn(123))
	fmt.Println(fn("screams"))
	// Output: 10
	// 20
	// 30
}

func ExampleImplementerOf() {
	var fn = func(s *mockMatcher) int {
		return 0
	}

	stub := mocka.Function(t, &fn, 10)
	defer stub.Restore()

	stub.WithArgs(match.ImplementerOf(new(match.SupportedKindsMatcher))).Return(20)

	fmt.Println(fn(nil))
	fmt.Println(fn(&mockMatcher{}))
	// Output: 10
	// 20
}

type mockMatcher struct {
}

// SupportedKinds returns the supported kinds for the matcher
func (mockMatcher) SupportedKinds() map[reflect.Kind]struct{} {
	return nil
}

// Match return true is the match was successful; otherwise false
func (mockMatcher) Match(interface{}) bool {
	return false
}
