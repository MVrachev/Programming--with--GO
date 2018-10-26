package main

import "fmt"

func Repeater(s, sep string) func(int) string {

	return func(count int) (result string) {
		for i := 0; i < count; i++ {
			result += s
			if i != count-1 {
				result += sep
			}
		}
		return result
	}
}

func Generator(gen func(int) int, initial int) func() int {

	return func() int {
		var old int = initial
		initial = gen(initial)
		return old
	}
}

func MapReducer(mapper func(int) int, reducer func(int, int) int, initial int) func(...int) int {

	return func(args ...int) int {
		result := initial
		for _, v := range args {
			mapVal := mapper(v)
			result = reducer(result, mapVal)
		}
		return result
	}
}

// func printArguments(args ...int) {
// 	var i int = 0
// 	for _, v := range args {
// 		fmt.Println("Argument number %d is: %d", i, v)
// 		i++
// 	}
// }
func main() {

	fmt.Println(Repeater("foo", ":")(-1))

	fmt.Println(Repeater("foo", ":")(3))

	counter := Generator(
		func(v int) int { return v + 1 },
		0)

	power := Generator(
		func(v int) int { return v * v },
		2,
	)

	fmt.Println("-------------------------")

	fmt.Println(counter()) // 0
	fmt.Println(counter()) // 1
	fmt.Println(power())   // 2
	fmt.Println(power())   // 4
	fmt.Println(counter()) // 2
	fmt.Println(power())   // 16
	fmt.Println(counter()) // 3
	fmt.Println(power())   // 256

	fmt.Println("-------------------------")

	powerSum := MapReducer(
		func(v int) int { return v * v },
		func(a, v int) int { return a + v },
		0,
	)

	fmt.Println(powerSum(1, 2, 3, 4)) // 30

	tripleSum := MapReducer(
		func(v int) int { return v * v * v },
		func(a, v int) int { return a + v },
		1000,
	)
	fmt.Println(tripleSum(1, 2, 3, 4))
	fmt.Println(tripleSum())

	fmt.Println("-------------------------")

	// printArguments(1, 2, 3, 4)

	idDivision := MapReducer(
		func(v int) int { return v },
		func(a, v int) int { return a / v },
		8,
	)

	fmt.Println(idDivision(4, 2)) // 1
	fmt.Println(idDivision())     // 8

}
