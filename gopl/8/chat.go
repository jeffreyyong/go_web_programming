package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// broadcaster
type client chan<- string // an outgoing message channel

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

// The broadcaster listens on the global entering and leaving channels for announcements of
// arriving and departing clients. When it receives one of these events, it updates the clients
// set, and if the event was a departure, it closes the client's outgoing message channel

// The broadcaster also listens for events on the global messages channel, to which each client sends
// all its incoming messages. When the broadcaster receives one of these events, it broadcasts the
// message to every connected client.

func broadcaster() {
	// The only information recorded about each client is the identity of its outgoing message channel
	clients := make(map[client]bool) // all incoming client messages
	for {
		select {
		case msg := <-messages:
			// broadcast incoming messages to all clients' outgoing message channels.
			for cli := range clients {
				fmt.Printf("\n iterated clients range: %v", cli)
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

// per-client goroutines. handleConn creates a new outgoing message channel
// for its client and announces the arrival of this client to the broadcaster
// over the entering channel

// It reads every line of text from the client, sending each line to the broadcaster
// over the global incoming message channel, prefixing each message with the identity
// of its sender. Once there is nothing more to read from the client, handleConn
// announces the departure of the client over the leaving channel and closes the connection

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
