package main

//import "../../internal/storage"
import (
	"log"
	"net/http"

	"github.com/rebus2015/praktikum-devops/internal/handlers"
	"github.com/rebus2015/praktikum-devops/internal/storage"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Panic(err)
	}

	storage := storage.CreateRepository()
	r := handlers.NewRouter(storage)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
