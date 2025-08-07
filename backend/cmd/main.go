package main

import (
	"log"
	"net"

	"github.com/delta-team-uz/grpc-web-todo/storage"
	"github.com/delta-team-uz/grpc-web-todo/todo_service_grpc"
	"google.golang.org/grpc"
)

func main() {
	todoService := storage.NewTodoService()
	lis, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatal("Error while listening: %w", err)
	}

	grpcServer := grpc.NewServer()

	todo_service_grpc.RegisterTodoServiceServer(grpcServer, todoService)

	log.Println("gRPC server listening on :5050")
	grpcServer.Serve(lis)
}
