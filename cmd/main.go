package main

import (
	"log"

	"github.com/evgeniyPP/go-todos-api"
	"github.com/evgeniyPP/go-todos-api/pkg/handler"
)

func main() {
	handlers := new(handler.Handler)
	srv := new(todos.Server)
	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server. Error: %s", err.Error())
	}
}
