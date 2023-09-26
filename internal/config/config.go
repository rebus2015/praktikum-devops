// Package config выполняет функцию параметризации сервиса сбора метрик
// Поддерживает задание параметров запуска через переменные окружения и параметры командной строки.
package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env"
)

// Config хранит получныые занчеия конфигурации.
type Config struct {
	ServerAddress    string        `env:"ADDRESS"`
	StoreInterval    time.Duration `env:"STORE_INTERVAL"` // 0 - синхронная запись
	StoreFile        string        `env:"STORE_FILE"`     // пустое значние отключает запись на диск
	Restore          bool          `env:"RESTORE"`        // загружать начальные значениея из файла
	Key              string        `env:"KEY"`            // Ключ для создания подписи сообщения
	ConnectionString string        `env:"DATABASE_DSN"`   // Cтрока подключения к БД
	CryptoKeyFile    string        `env:"CRYPTO_KEY"`     // путь к файлу с приватным ключом
	CryptoKey        *rsa.PrivateKey
}

// GetConfig считывает значения параметров запуска и возвращает структуру.
func GetConfig() (*Config, error) {
	conf := Config{}

	flag.StringVar(&conf.ServerAddress, "a", "127.0.0.1:8080", "Server address")
	flag.DurationVar(&conf.StoreInterval, "i", time.Second*30, "Metrics save to file interval")
	flag.StringVar(&conf.StoreFile, "f", "", "Metrics repository file path")
	flag.BoolVar(&conf.Restore, "r", false, "Restore metric values from file before start")
	flag.StringVar(&conf.Key, "k", "", "Key to sign up data with SHA256 algorythm")
	flag.StringVar(&conf.ConnectionString, "d", "",
		"Database connection string(PostgreSql)") // postgresql://pguser:pgpwd@localhost:5432/devops?sslmode=disable
	flag.StringVar(&conf.CryptoKeyFile, "crypto-key", "/Users/mak/go/praktikum-devops/keys/privateKey.pem", "Public Key file address")
	flag.Parse()

	err := env.Parse(&conf)
	if err != nil {
		return nil, fmt.Errorf("error reading agent  config: %w", err)
	}
	if err = conf.getCryptoKey(); err != nil {
		return nil, fmt.Errorf("error reading agent config, failed to get CryptoKey: %w", err)
	}

	return &conf, err
}

func (c *Config) getCryptoKey() error {
	if c.CryptoKeyFile == "" {
		return nil
	}
	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if _, err := os.Stat(c.CryptoKeyFile); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error reading agent config: %w", err)
	}
	filename := c.CryptoKeyFile
	//1. Read the private key information and put it in the data variable
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	stat, _ := file.Stat() //Get file attribute information
	data := make([]byte, stat.Size())
	file.Read(data)
	file.Close()
	//2. Decode the resulting string pem
	block, _ := pem.Decode(data)

	privateKey, err3 := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err3 != nil {
		return fmt.Errorf("error reading agent  config: %w", err3)
	}
	c.CryptoKey = privateKey
	return nil
}
