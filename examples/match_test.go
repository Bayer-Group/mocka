package examples

import (
	"fmt"
	"log"

	"github.com/MonsantoCo/mocka"
	"github.com/MonsantoCo/mocka/match"
)

func ExampleAnything() {
	var fn = func(str []string, n int) int {
		return len(str) + n
	}

	stub, err := mocka.Function(&fn, 20)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer stub.Restore()

	if err = stub.WithArgs(match.Anything(), 2).Return(10); err != nil {
		log.Fatal(err.Error())
	}

	if err = stub.WithArgs([]string{"123", "456"}, 2).Return(5); err != nil {
		log.Fatal(err.Error())
	}

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

	stub, err := mocka.Function(&fn, 20)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer stub.Restore()

	if err = stub.WithArgs(match.AnythingButNil(), 2).Return(10); err != nil {
		log.Fatal(err.Error())
	}

	if err = stub.WithArgs([]string{"123", "456"}, 2).Return(5); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(fn([]string{"hello"}, 5))
	fmt.Println(fn([]string{"mocka"}, 2))
	fmt.Println(fn(nil, 2))
	fmt.Println(fn([]string{"123", "456"}, 2))
	// Output: 20
	// 10
	// 20
	// 5
}
