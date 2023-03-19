package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddress string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"` // 0 - синхронная запись
	StoreFile     string        `env:"STORE_FILE"`     // пустое значние отключает запись на диск
	Restore       bool          `env:"RESTORE"`        // загружать начальные значениея из файла
}

func GetConfig() (*Config, error) {
	conf := Config{}

	flag.StringVar(&conf.ServerAddress, "a", "127.0.0.1:8080", "Server address")
	flag.DurationVar(&conf.StoreInterval, "i", time.Second*300, "Metrics save to file interval")
	flag.StringVar(&conf.StoreFile, "f", "/tmp/devops-metrics-db.json", "Metrics repository file path")
	flag.BoolVar(&conf.Restore, "r", true, "Restore metric values from file before start")

	flag.Parse()

	err := env.Parse(&conf)

	return &conf, err
}
