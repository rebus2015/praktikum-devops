package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/rebus2015/praktikum-devops/internal/config"
	"github.com/rebus2015/praktikum-devops/internal/handlers"
	"github.com/rebus2015/praktikum-devops/internal/storage"
	"github.com/rebus2015/praktikum-devops/internal/storage/dbstorage"
	"github.com/rebus2015/praktikum-devops/internal/storage/filestorage"
	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Panicf("Error reading configuration from env variables: %v", err)
		return
	}
	log.Printf("server started on %v with key: '%v'", cfg.ServerAddress, cfg.Key)

	fs := filestorage.NewStorage(cfg)
	var ms = new(memstorage.MemStorage)
	if cfg.Restore && cfg.StoreFile != "" {
		ms = fs.Restore(cfg.StoreFile)
	} else {
		ms = memstorage.NewStorage()
	}
	if cfg.StoreInterval != 0 {
		go fs.SaveTicker(cfg.StoreInterval, ms)
	}

	var storage storage.Repository = storage.NewRepositoryWrapper(*ms, *fs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sqlDBStorage, err := dbstorage.NewPostgreSQLStorage(ctx, cfg.ConnectionString)
	if err != nil {
		log.Printf("Error creating dbStorage: %v", err)
		log.Panicf("Error creating dbStorage: %v", err)
		return
	}
	log.Printf("Created dbStorage: %v", cfg.ConnectionString)
	// defer sqlDBStorage.Close()

	r := handlers.NewRouter(&storage, sqlDBStorage, *cfg)
	srv := &http.Server{
		Addr:         cfg.ServerAddress,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      r,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Printf("server exited with %v", err)
	}

}
