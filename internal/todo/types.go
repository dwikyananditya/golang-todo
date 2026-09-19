package todo

import "gorm.io/gorm"

type Model struct {
	gorm.Model
	Title       string
	Description string
	IsDone      bool
}

type Todo struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	IsDone      bool   `json:"is_done"`
}

func (m Model) Todo() Todo {
	return Todo{
		Title:       m.Title,
		Description: m.Description,
		IsDone:      m.IsDone,
	}
}
