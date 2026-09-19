package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"todo/internal/todo"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	todos todo.Repository
}

func New(todos todo.Repository) *Handler {
	return &Handler{todos: todos}
}

type StructValidator struct {
	validate *validator.Validate
}

func NewStructValidator() *StructValidator {
	return &StructValidator{validate: validator.New()}
}

func (v *StructValidator) Validate(value any) error {
	return formatValidationErrors(v.validate.Struct(value))
}

func formatValidationErrors(err error) error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err
	}

	messages := make([]string, len(validationErrors))
	for i, validationError := range validationErrors {
		messages[i] = fmt.Sprintf("field '%s' failed on '%s'", validationError.Field(), validationError.Tag())
	}
	return errors.New(strings.Join(messages, "; "))
}

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func respondOK(c fiber.Ctx, data any) error {
	return c.JSON(Response{Data: data})
}

func respondError(c fiber.Ctx, status int, message string, err error) error {
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}
	return c.Status(status).JSON(Response{Message: message, Error: errorMessage})
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Patch("/:id", h.update)
	app.Delete("/:id", h.delete)
	app.Get("/", h.get)
	app.Post("/", h.create)
}

func (h *Handler) update(c fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id", err)
	}

	request := new(todo.Todo)
	if err := c.Bind().Body(request); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body", err)
	}

	rows, err := h.todos.Update(context.Background(), id, *request)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to update todo", err)
	}
	if rows == 0 {
		return respondError(c, fiber.StatusNotFound, "todo not found", nil)
	}
	return respondOK(c, request)
}

func (h *Handler) delete(c fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id", err)
	}

	rows, err := h.todos.Delete(context.Background(), id)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to delete todo", err)
	}
	if rows == 0 {
		return respondError(c, fiber.StatusNotFound, "todo not found", nil)
	}
	return respondOK(c, fiber.Map{"id": id})
}

func (h *Handler) get(c fiber.Ctx) error {
	item, err := h.todos.First(context.Background())
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "todo not found", err)
	}
	return respondOK(c, item)
}

func (h *Handler) create(c fiber.Ctx) error {
	request := new(todo.Todo)
	if err := c.Bind().Body(request); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body", err)
	}

	if err := h.todos.Create(context.Background(), *request); err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create todo", err)
	}
	return respondOK(c, request)
}

func parseID(c fiber.Ctx) (int, error) {
	return strconv.Atoi(c.Params("id"))
}
