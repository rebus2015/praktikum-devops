package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddress    string        `env:"ADDRESS"`
	StoreInterval    time.Duration `env:"STORE_INTERVAL"` // 0 - синхронная запись
	StoreFile        string        `env:"STORE_FILE"`     // пустое значние отключает запись на диск
	Restore          bool          `env:"RESTORE"`        // загружать начальные значениея из файла
	Key              string        `env:"KEY"`            // Ключ для создания подписи сообщения
	ConnectionString string        `env:"DATABASE_DSN"`   // Cтрока подключения к БД
}

func GetConfig() (*Config, error) {
	conf := Config{}

	flag.StringVar(&conf.ServerAddress, "a", "127.0.0.1:8080", "Server address")
	flag.DurationVar(&conf.StoreInterval, "i", time.Second*20, "Metrics save to file interval")
	flag.StringVar(&conf.StoreFile, "f", "/tmp/devops-metrics-db.json", "Metrics repository file path")
	flag.BoolVar(&conf.Restore, "r", true, "Restore metric values from file before start")
	flag.StringVar(&conf.Key, "k", "", "Key to sign up data with SHA256 algorythm")
	flag.StringVar(&conf.ConnectionString, "d", "",
		"Database connection string(PostgreSql)") // postgresql://pguser:pgpwd@localhost:5432/devops?sslmode=disable

	flag.Parse()

	err := env.Parse(&conf)

	return &conf, err
}
