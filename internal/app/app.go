package app

import (
	"todo/internal/handler"
	"todo/internal/pb"
	"todo/internal/todo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func New(todoRepository todo.Repository) *grpc.Server {
	server := grpc.NewServer()
	pb.RegisterTodoServiceServer(server, handler.New(todoRepository))
	reflection.Register(server)
	return server
}
