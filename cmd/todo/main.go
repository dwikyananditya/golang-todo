package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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
	server, err := app.New(todoRepository)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening on :8080")
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("shutting down")
		server.GracefulStop()
	}()
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
