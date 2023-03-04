package main

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type config struct {
	ServerAddress  string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInternal time.Duration `env:"REPORT_INTERVAL" envDefault:"5s"`
	PollInterval   time.Duration `env:"POLL_INREVAL" envDefault:"2s"`
}

func getConfig() (*config, error) {
	var conf config
	err := env.Parse(&conf)
	return &conf, err
}
