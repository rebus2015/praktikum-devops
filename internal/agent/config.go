// Package agent реализует агент сбора метрик
package agent

import (
	"flag"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddress  string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"PUSH_TIMEOUT"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
	RateLimit      int           `env:"RATE_LIMIT"` // Количество одновременно исходящих запросов на сервер
}

func GetConfig() (*Config, error) {
	conf := Config{}
	flag.StringVar(&conf.ServerAddress, "a", "127.0.0.1:8080", "Server address")
	flag.DurationVar(&conf.ReportInterval, "r", time.Second*11, "Interval before push metrics to server")
	flag.DurationVar(&conf.PollInterval, "p", time.Second*5, "Interval between metrics reads from runtime")
	flag.StringVar(&conf.Key, "k", "", "Key to sign up data with SHA256 algorythm")
	flag.IntVar(&conf.RateLimit, "l", 12, "Workers count")
	flag.Parse()
	err := env.Parse(&conf)

	return &conf, err
}
