package main

import (
	"log"
	"net/rpc"
)

// Client dials to the server and get the remote function executed
// The only way to get data back is to pass the reply pointer object along with request

type Args struct {
}

func main() {
	var reply int64
	args := Args{}
	// Do a DialHTTP to connect to the RPC server, which is running on the localhost 1234
	client, err := rpc.DialHTTP("tcp", "localhost"+":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Call the Remote function with the Name:Function format with args and reply with the pointer object.
	err = client.Call("TimeServer.GiveServerTime", args, &reply)
	if err != nil {
		log.Fatal("arith err:", err)
	}
	// Get the data collected into the reply object
	log.Printf("%d", reply)
}
