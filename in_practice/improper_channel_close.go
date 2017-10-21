package main

import (
	"fmt"
	"time"
)

func main() {
	msg := make(chan string)
	until := time.After(5 * time.Second)

	// Starts a send goroutine with a sending channel
	go send(msg)

	// Loops over a select that watches for messages from send, or for a time-out
	for {
		select {
		// If a message arrives frmo send, prints it
		case m := <-msg:
			fmt.Println(m)
		case <-until:
			// When the time-out occurs, shut things down. Pause to ensure
			// that the failure can be seen before the main goroutine exits.
			close(msg)
			time.Sleep(500 * time.Millisecond)
			return
		}
	}
}

func send(ch chan string) {
	// Sends "Hello" to the channel every half-second
	for {
		ch <- "hello"
		time.Sleep(500 * time.Millisecond)
	}
}

// At the end, the program panics because main closes the msg channel while send is still sending
// messages to it. A send on a closed channel panics. The close function should be closed only by a sender,
// in general it should be done with some protective guards around it.
