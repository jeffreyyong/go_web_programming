package main

import "fmt"

func counter(out chan<- int) {
	for x := 0; x < 100000; x++ {
		out <- x
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for v := range in {
		out <- v * v
	}
	close(out)
}

func printer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// This implicitly converts naturals, a value of type chan int, to the type of the parameter, chan<- int
	go counter(naturals)
	go squarer(squares, naturals)
	// This does a similar implicit conversion to <-chan int
	printer(squares)
}

// Conversions fmor bidirectional to unidirectional channel types are permitted in any assignment
// There is no goibg back, once a value of a unidirectional type such as chan<- int, there is
// no way to obtain from it a value of type chan int that regeres to the same channel data structure.
