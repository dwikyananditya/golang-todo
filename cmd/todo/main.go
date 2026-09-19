package main

import (
	"log"

	"todo/internal/app"
	"todo/internal/todo"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	todoRepository := todo.NewGormRepository(db)
	if err := todoRepository.Migrate(); err != nil {
		log.Fatal(err)
	}

	application := app.New(todoRepository)
	log.Fatal(application.Listen(":8080"))
}
