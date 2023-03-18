package main

import (
	"log"
	"net/http"

	"github.com/rebus2015/praktikum-devops/internal/config"
	"github.com/rebus2015/praktikum-devops/internal/handlers"
	"github.com/rebus2015/praktikum-devops/internal/storage"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Panicf("Error reading configuration from env variables: %v", err)
		return
	}
	log.Printf("server started on %v", cfg.ServerAddress)
	storage := storage.Create(cfg)
	r := handlers.NewRouter(&storage)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
