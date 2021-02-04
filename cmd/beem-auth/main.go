package main

import (
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"

	"beem-auth/internal/pb"
	"beem-auth/internal/pkg/database"
	service "beem-auth/internal/pkg/service"
)

const (
	port = ":5051"
)

func main() {
	log.Println("Running GRPC with version", grpc.Version)

	db, err := database.Connect("localhost", "5432", "postgres", "postgres", "beem_auth")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Listen to tcp port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		// Interceptors are executed from left to right (or top to bottom as seen here)
		grpc_middleware.WithUnaryServerChain(),
	)

	pb.RegisterAccountServiceServer(grpcServer, service.NewAccountController(db))

	log.Printf("started...")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
