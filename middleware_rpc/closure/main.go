package main

import "fmt"

func main() {
	numGenerator := generator()
	for i := 0; i < 5; i++ {
		fmt.Print(numGenerator(), "\t")
	}
}

// This function returns another function
// A closure function called generator and calling it to get a new number.
// A generator pattern generates a new item each time, based on given conditions
// THe inner function getting returned is an anonymous function with no arguments and one return type of integer
// The variable i that is defined inside the outer function is available to the anonymous function
// The function signature of the outer function should exactly match the anonymous function's signature
//

func generator() func() int {
	var i = 0
	return func() int {
		i++
		return i
	}
}
