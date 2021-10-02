package main

import (
	"log"

	"github.com/evgeniyPP/go-todos-api"
	"github.com/evgeniyPP/go-todos-api/pkg/handler"
	"github.com/evgeniyPP/go-todos-api/pkg/repository"
	"github.com/evgeniyPP/go-todos-api/pkg/service"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initilizing configs. Error: %s", err.Error())
	}

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todos.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server. Error: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
