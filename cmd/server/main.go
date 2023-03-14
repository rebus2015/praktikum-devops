package main

//import "../../internal/storage"
import (
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/rebus2015/praktikum-devops/internal/handlers"
	"github.com/rebus2015/praktikum-devops/internal/storage"
)

type config struct {
	ServerAddress string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

func getConfig() (*config, error) {
	conf := config{}
	err := env.Parse(&conf)

	return &conf, err
}
func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Panicf("Error reading configuration from env variables: %v", err)
		return
	}
	log.Printf("server started on %v", cfg.ServerAddress)
	storage := storage.CreateRepository()
	r := handlers.NewRouter(storage)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
