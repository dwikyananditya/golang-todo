package todo

import "gorm.io/gorm"

type Model struct {
	gorm.Model
	Title       string
	Description string
	IsDone      bool
}

type Todo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsDone      bool   `json:"is_done"`
}

func (m Model) Todo() Todo {
	return Todo{
		Title:       m.Title,
		Description: m.Description,
		IsDone:      m.IsDone,
	}
}
