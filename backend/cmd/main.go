package main

import (
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/delta-team-uz/grpc-web-todo/storage"
	"github.com/delta-team-uz/grpc-web-todo/todo_service_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

func main() {
	todoService := storage.NewTodoService()
	lis, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatal("Error while listening: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    time.Minute,
			Timeout: 5,
		}),
		grpc.MaxRecvMsgSize(10*1024*1024), // 10MB
		grpc.MaxSendMsgSize(10*1024*1024), // 10MB
	)

	reflection.Register(grpcServer)
	todo_service_grpc.RegisterTodoServiceServer(grpcServer, todoService)
	slog.Info("grpc server running port :5050 ")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Error while listening: %w", err)
	}

}
