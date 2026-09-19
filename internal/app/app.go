package app

import (
	"todo/internal/handler"
	"todo/internal/todo"

	"github.com/gofiber/fiber/v3"
)

func New(todoRepository todo.Repository) *fiber.App {
	application := fiber.New(fiber.Config{
		StructValidator: handler.NewStructValidator(),
	})

	handler.New(todoRepository).RegisterRoutes(application)
	return application
}
