package main

import (
	"fmt"
	"time"
)

func main() {
	// Creates a buffered channel with one space
	lock := make(chan bool, 1)
	// Starts up to six goroutines sharing the lock channel
	for i := 1; i < 7; i++ {
		go worker(i, lock)
	}
	time.Sleep(10 * time.Second)
}

func worker(id int, lock chan bool) {
	fmt.Printf("%d wants the lock\n", id)
	// A worker acquires the lock by sending it a message.
	// The first worker to his this will get the one space, and thus own the lock. The rest will block
	lock <- true
	fmt.Printf("%d has the lock\n", id)
	time.Sleep(500 * time.Millisecond)
	// The space between "lock <- true" and "<- lock" is locked.
	fmt.Printf("%d is releasing the lock\n", id)
	// Releases the lock by reading a value, which then opens the space on the buffer again
	// so that the next function can lock it
	<-lock
}
