package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// ...create abort channel...
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()

	fmt.Println("Commencing countdown. Press return to abort.")

	// The select statement below waits until the first of two events arrives, either an abort event or the event indicating
	// that 10 seconds have elapsed. If 10 seconds go by with no abort, the launch proceeds
	select {
	// time.After immediately returns a channel, and starts a new goroutine that sends a single value on that channel after the specified time.
	case <-time.After(10 * time.Second):
	// Do nothing
	case <-abort:
		fmt.Println("Launch aborted!")
		return
	}
	launch()
}

func launch() {
	fmt.Println("Lift off!")
}
