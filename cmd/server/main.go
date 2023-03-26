package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rebus2015/praktikum-devops/internal/config"
	"github.com/rebus2015/praktikum-devops/internal/handlers"
	"github.com/rebus2015/praktikum-devops/internal/storage"
	"github.com/rebus2015/praktikum-devops/internal/storage/dbstorage"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Panicf("Error reading configuration from env variables: %v", err)
		return
	}
	log.Printf("server started on %v with key: '%v'", cfg.ServerAddress, cfg.Key)
	storage := storage.Create(cfg)

	sqlDBStorage, err := dbstorage.NewPostgreSQLStorage(cfg.ConnectionString)
	if err != nil {
		log.Printf("Error creating dbStorage: %v", err)
		log.Panicf("Error creating dbStorage: %v", err)
		return
	}
	// defer sqlDBStorage.Close()

	r := handlers.NewRouter(&storage, sqlDBStorage, *cfg)
	srv := &http.Server{
		Addr:         cfg.ServerAddress,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      r,
	}
	log.Fatal(srv.ListenAndServe())
}
