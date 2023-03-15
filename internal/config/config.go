package config

import (
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddress string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"30s"`                    //0 - синхронная запись
	StoreFile     string	    `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"` //пустое значние отключает запись на диск
	Restore       bool          `env:"RESTORE" envDefault:"true"`                           //загружать начальные значениея из файла
}

func GetConfig() (*Config, error) {
	conf := Config{}
	err := env.Parse(&conf)

	return &conf, err
}
