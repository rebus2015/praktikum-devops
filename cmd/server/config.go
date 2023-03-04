package main
import (
	"github.com/caarlos0/env/v6"
)

type config struct {
	ServerAddress  string	`env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

func getConfig() (*config, error) {
	var conf config
	err := env.Parse(&conf)
	return &conf, err
}