package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "go_web_programming/protobuf_grpc/grpc/datafiles"
)

const (
	port = ":50051"
)

// server is used to create MoneyTransactionServer
type server struct{}

// MakeTransaction implements MoneyTransactionServer.MakeTransaction
// Context is used to create a context variable, which lives throughout an RPC request's lifetime
func (s *server) MakeTransaction(ctx context.Context, in *pb.TransactionRequest) (*pb.TransactionResponse, error) {
	log.Printf("Got request for money Transfer....")
	log.Printf("Amount: %f, from A/c:%s, to A/c:%s", in.Amount, in.From, in.To)
	// Do database logic herer
	return &pb.TransactionResponse{Confirmation: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMoneyTransactionServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
