package main

import (
	"fmt"
	"time"
)

func main() {
	msg := make(chan string)
	// An additional Boolean channel that indicates when it's finished
	done := make(chan bool)
	until := time.After(5 * time.Second)

	// Passes two channels into send
	go send(msg, done)

	for {
		select {
		case m := <-msg:
			fmt.Println(m)
		case <-until:
			// When timeout, lets "send" know the process is done
			done <- true
			time.Sleep(500 * time.Millisecond)
			return
		}
	}
}

// ch is a receiving channel, while done is a sending channel
func send(ch chan<- string, done <-chan bool) {
	for {
		select {
		// When done has a message, shuts things down
		case <-done:
			println("Done")
			close(ch)
			return
		default:
			ch <- "hello"
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// In this pattern, have one goroutine whose primary task is to receive messages,
// And another whose job is to send messages.
// If the receiver hits a stopping condition, it must let the sender know.

// The main function is the one that knows when to stop processing, but it's also the receiver.
// The receiver shouldn't ever close a receiving channel.
// Instead, it sends a message on the done channel indicating that it's done with the work
// The send function knows when it receives the message "done" that it can close the channel and return.
