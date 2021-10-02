package main

import (
	"log"

	"github.com/evgeniyPP/go-todos-api"
	"github.com/evgeniyPP/go-todos-api/pkg/handler"
	"github.com/evgeniyPP/go-todos-api/pkg/repository"
	"github.com/evgeniyPP/go-todos-api/pkg/service"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todos.Server)
	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server. Error: %s", err.Error())
	}
}
