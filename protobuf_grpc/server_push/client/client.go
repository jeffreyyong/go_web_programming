package main

import (
	"context"
	"io"
	"log"

	pb "go_web_programming/protobuf_grpc/server_push/datafiles"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

// ReceiveStream listens to the stream contents and use them
// It uses the first argument to create a stream and start listening to it. Whenever the server
// exhausts all the messages, the client will stop listening and terminate. Then, an io.EOF error will
// be returned if the client tries to receive messages.
// The second argument, TransactionRequest is used to send the request to the server for the first time.
func ReceiveStream(client pb.MoneyTransactionClient, request *pb.TransactionRequest) {
	log.Println("Started listening ot the server stream!")
	stream, err := client.MakeTransaction(context.Background(), request)
	if err != nil {
		log.Fatalf("%v.MakeTransaction(_) = _, %v", client, err)
	}
	// Listen to teh stream of messages
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// If there are no more responses, get out of loop
			break
		}
		if err != nil {
			log.Fatalf("%v.MakeTransaction(_) = _, %v", client, err)
		}
		log.Printf("Status: %v, Operation: %v", response.Status, response.Description)
	}
}

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMoneyTransactionClient(conn)

	// Prepare data. Get this from clients like front-end or app
	from := "1234"
	to := "5678"
	amount := float32(1250.75)

	// Contact the server and print out its response.
	ReceiveStream(client, &pb.TransactionRequest{From: from, To: to, Amount: amount})
}

// The client stays alive until all the streaming messages are sent back. The server can handle any number of clients
// at a given time. Every client request is consiered as an individual entity.
