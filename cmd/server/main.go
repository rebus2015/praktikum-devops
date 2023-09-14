package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // #nosec
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/rebus2015/praktikum-devops/internal/config"
	"github.com/rebus2015/praktikum-devops/internal/handlers"
	"github.com/rebus2015/praktikum-devops/internal/storage"
	"github.com/rebus2015/praktikum-devops/internal/storage/dbstorage"
	"github.com/rebus2015/praktikum-devops/internal/storage/filestorage"
	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	fmt.Printf("version=%s, date=%s\n, version=%s", buildVersion, buildDate, buildVersion)
	cfg, err := config.GetConfig()
	if err != nil {
		log.Panicf("Error reading configuration from env variables: %v", err)
		return
	}
	log.Printf("server started on %v with \n key: '%v', \n store.Interval:%v,\n restore: %v ",
		cfg.ServerAddress, cfg.Key, cfg.StoreInterval, cfg.Restore)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var fs storage.SecondaryStorage
	var sqlDBStorage dbstorage.SQLStorage
	switch {
	case cfg.ConnectionString != "":
		db, dberr := dbstorage.NewStorage(ctx, cfg.ConnectionString, true)
		if dberr != nil {
			log.Panicf("Error creating dbstorage %v", err)
		}
		fs = storage.SecondaryStorage(db)
		sqlDBStorage = dbstorage.SQLStorage(db)
		defer sqlDBStorage.Close()
	default:
		if cfg.StoreFile != "" {
			fs = filestorage.NewStorage(cfg)
		}
	}

	ms := memstorage.NewStorage()
	if fs != nil {
		if cfg.Restore {
			ms, err = fs.Restore()
			if err != nil {
				log.Panicf("Error restoring data %v", err)
			}
		}
		if cfg.StoreInterval != 0 {
			go fs.SaveTicker(cfg.StoreInterval, ms)
		}
	}

	storage := storage.NewRepositoryWrapper(ms, fs)
	defer cancel()
	if err != nil {
		log.Printf("Error creating NewRepositoryWrapper: %v", err)
	}
	log.Printf("Created NewRepositoryWrapper: %v", cfg.ConnectionString)

	r := handlers.NewRouter(storage, sqlDBStorage, *cfg)
	srv := &http.Server{
		Addr:         cfg.ServerAddress,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      r,
	}
	err = srv.ListenAndServe()
	log.Printf("server started \n address:%v \n database:%v,\n restore interval: %v ",
		cfg.ServerAddress, cfg.ConnectionString, cfg.StoreInterval)

	if err != nil {
		log.Printf("server exited with %v", err)
	}
}
