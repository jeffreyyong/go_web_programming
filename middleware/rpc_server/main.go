package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

// Holds informatoin about arguments passed from client to the server (both RPC).
type Args struct{}

type TimeServer int64

// Function that will be called by the client, and the current server time is returned
// GiveServerTime takes the Args object as the first argument and a reply pointer object
// It sets the reply pointer object but does not return anything except an error
// The Args struct here has no fields because this server is not expecting the client to send any arguments
func (t *TimeServer) GiveServerTime(args *Args, reply *int64) error {
	*reply = time.Now().Unix()
	return nil
}

func main() {
	timeserver := new(TimeServer)
	// Create a TimeServer number to register with the rpc.Register
	// The server wishes to export an object of type TimeServer int64
	rpc.Register(timeserver)
	// Registers an HTTP handler for RPC messages to DefaultServer.
	rpc.HandleHTTP()
	// Start a TCP server that listens on port 1234
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	// Serve the running progrgamme
	http.Serve(l, nil)
}
