package main

import (
	"log"
	"net"
	"os"

	"github.com/bobadojo/go/pkg/stores/v1/storespb"
	"google.golang.org/grpc"
)

func main() {
	storesServer, err := NewStoresServer()
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on port %s", port)
	grpcServer := grpc.NewServer()
	storespb.RegisterStoresServer(grpcServer, storesServer)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
