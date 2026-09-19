package app

import (
	"buf.build/go/protovalidate"
	"todo/internal/handler"
	"todo/internal/pb"
	"todo/internal/todo"

	pbmw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func New(todoRepository todo.Repository) (*grpc.Server, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		pbmw.UnaryServerInterceptor(validator),
	))
	pb.RegisterTodoServiceServer(server, handler.New(todoRepository))
	reflection.Register(server)
	return server, nil
}
