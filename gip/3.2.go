package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Outside a goroutine")
	go func() {
		fmt.Println("Inside a goroutine")
	}()
	fmt.Println("Outside again.")

	// Yields to the scheduler

	runtime.Gosched()
}

// Goroutines run concurrently but not necessarily in parallel. When a goroutine is scheduled to run by calling go fun
// the Go runtime executes the function as soon as it can.

// runtime.Gosched() indicates to the Go runtime that at this point, it yields to the scheduler.
// If the scheduler has other tasks queued up (other goroutines), it may then run one or more of them before coming back to this function.
