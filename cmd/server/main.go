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
	ctx, cancel := context.WithCancel(context.Background())

	var fs storage.SecondaryStorage
	switch {
	case cfg.ConnectionString != "":
		fs, err = dbstorage.NewStorage(ctx, cfg.ConnectionString, cfg.StoreInterval == 0)
		if err != nil {
			log.Panicf("Error creating dbstorage %v", err)
		}
	default:
		if cfg.StoreFile != "" {
			fs = filestorage.NewStorage(cfg)
		}
	}

	ms := memstorage.NewStorage()
	if cfg.Restore && fs != nil {
		ms, err = fs.Restore()
		if err != nil {
			log.Panicf("Error restoring data %v", err)
		}
	}
	if cfg.StoreInterval != 0 {
		go fs.SaveTicker(cfg.StoreInterval, ms)
	}

	var storage storage.Repository = storage.NewRepositoryWrapper(*ms, fs)

	defer cancel()
	sqlDBStorage, err := dbstorage.NewStorage(ctx, cfg.ConnectionString, false)
	if err != nil {
		log.Printf("Error creating dbStorage: %v", err)
	}
	log.Printf("Created dbStorage: %v", cfg.ConnectionString)
	defer sqlDBStorage.Close()

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
