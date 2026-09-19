package main

import (
	"log"
	"net"

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

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening on :8080")
	log.Fatal(app.New(todoRepository).Serve(lis))
}
