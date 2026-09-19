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

func (h *Handler) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.GetTodoResponse, error) {
	item, err := h.todos.First(ctx, req.GetOrder() == pb.SortOrder_SORT_ORDER_DESC)
	if err != nil {
		return nil, status.Error(codes.NotFound, "todo not found")
	}
	return &pb.GetTodoResponse{Todo: toPb(int64(item.ID), item.Todo())}, nil
}

func (h *Handler) ListTodos(ctx context.Context, req *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	items, err := h.todos.List(ctx, req.GetOrder() == pb.SortOrder_SORT_ORDER_DESC)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list todos")
	}
	out := make([]*pb.Todo, 0, len(items))
	for _, m := range items {
		out = append(out, toPb(int64(m.ID), m.Todo()))
	}
	return &pb.ListTodosResponse{Todos: out}, nil
}

func (h *Handler) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.UpdateTodoResponse, error) {
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

func toPb(id int64, item todo.Todo) *pb.Todo {
	return &pb.Todo{
		Id:          id,
		Title:       item.Title,
		Description: item.Description,
		IsDone:      item.IsDone,
	}
}
