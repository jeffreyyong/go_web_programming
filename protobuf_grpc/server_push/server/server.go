package main

import (
	"fmt"
	"log"
	"net"

	pb "go_web_programming/protobuf_grpc/server_push/datafiles"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port      = ":50051"
	noOfSteps = 3
)

// server is used to create MoneyTransactionServer
type server struct{}

// MakeTransaction implements MoneyTransactionServer.MakeTransaction
// Looping through the number of steps and performing the computation
// The server is simulationg the mock I/O or computation using the time.Sleep function
// This function sends a stream response from the server to the client.
func (s *server) MakeTransaction(in *pb.TransactionRequest, stream pb.MoneyTransaction_MakeTransactionServer) error {
	log.Printf("Got request for money transfer...")
	log.Printf("Amount: $%f, From A/c:%s, To A/c:%s", in.Amount, in.From, in.To)
	// Send streams here
	for i := 0; i < noOfSteps; i++ {
		// Simulating I/O or computation process with sleep..........
		// Usually this will be saving money trasnfer details in DB or talk to the third party API
		// time.Sleep(time.Second * 2)
		// Once task is done, send the successful message back to the client
		if err := stream.Send(&pb.TransactionResponse{Status: "good",
			Step:        int32(i),
			Description: fmt.Sprintf("Description of step %d", int32(i))}); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, "status", err)
		}
	}
	log.Printf("Successfully transfered amount $%v from %v to %v", in.Amount, in.From, in.To)
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	// Create a new gRPC server
	s := grpc.NewServer()
	// Register it with Proto service
	pb.RegisterMoneyTransactionServer(s, &server{})
	// Register reflection service on gRPC server
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
