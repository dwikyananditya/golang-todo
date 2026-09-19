package handler

import (
	"context"

	"todo/internal/pb"
	"todo/internal/todo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedTodoServiceServer
	todos todo.Repository
}

func New(todos todo.Repository) *Handler {
	return &Handler{todos: todos}
}

func (h *Handler) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	if err := require(req.Title, req.Description); err != nil {
		return nil, err
	}
	item := todo.Todo{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		IsDone:      req.GetIsDone(),
	}
	id, err := h.todos.Create(ctx, item)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create todo")
	}
	return &pb.CreateTodoResponse{Todo: toPb(id, item)}, nil
}

func (h *Handler) GetTodo(ctx context.Context, _ *pb.GetTodoRequest) (*pb.GetTodoResponse, error) {
	item, err := h.todos.First(ctx)
	if err != nil {
		return nil, status.Error(codes.NotFound, "todo not found")
	}
	return &pb.GetTodoResponse{Todo: toPb(int64(item.ID), item.Todo())}, nil
}

func (h *Handler) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.UpdateTodoResponse, error) {
	if err := require(req.Title, req.Description); err != nil {
		return nil, err
	}
	item := todo.Todo{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		IsDone:      req.GetIsDone(),
	}
	rows, err := h.todos.Update(ctx, int(req.GetId()), item)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update todo")
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, "todo not found")
	}
	return &pb.UpdateTodoResponse{Todo: toPb(req.GetId(), item)}, nil
}

func (h *Handler) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	rows, err := h.todos.Delete(ctx, int(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete todo")
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, "todo not found")
	}
	return &pb.DeleteTodoResponse{}, nil
}

func require(title, description string) error {
	if title == "" || description == "" {
		return status.Error(codes.InvalidArgument, "title and description are required")
	}
	return nil
}

func toPb(id int64, item todo.Todo) *pb.Todo {
	return &pb.Todo{
		Id:          id,
		Title:       item.Title,
		Description: item.Description,
		IsDone:      item.IsDone,
	}
}
