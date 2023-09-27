// Package config выполняет функцию параметризации сервиса сбора метрик
// Поддерживает задание параметров запуска через переменные окружения и параметры командной строки.
package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/caarlos0/env"
)

// Config хранит получныые занчеия конфигурации.
type Config struct {
	ServerAddress    string        `env:"ADDRESS" json:"address"`
	StoreInterval    time.Duration `env:"STORE_INTERVAL" json:"store_interval"` // 0 - синхронная запись
	StoreFile        string        `env:"STORE_FILE" json:"store_file"`         // пустое значние отключает запись на диск
	Restore          bool          `env:"RESTORE" json:"restore"`               // загружать начальные значениея из файла
	Key              string        `env:"KEY"`                                  // Ключ для создания подписи сообщения
	ConnectionString string        `env:"DATABASE_DSN" json:"database_dsn"`     // Cтрока подключения к БД
	CryptoKeyFile    string        `env:"CRYPTO_KEY" json:"crypto_key"`         // путь к файлу с приватным ключом
	confFile         string        `env:"CONFIG" json:"-"`
	CryptoKey        *rsa.PrivateKey
}

// GetConfig считывает значения параметров запуска и возвращает структуру.
func GetConfig() (*Config, error) {
	conf := Config{}
	flag.StringVar(&conf.confFile, "config", "", "Pass the conf.json path")
	flag.StringVar(&conf.confFile, "c", "", "Pass the conf.json path (shorthand)")
	flag.StringVar(&conf.ServerAddress, "a", "127.0.0.1:8080", "Server address")
	flag.DurationVar(&conf.StoreInterval, "i", time.Second*30, "Metrics save to file interval")
	flag.StringVar(&conf.StoreFile, "f", "", "Metrics repository file path")
	flag.BoolVar(&conf.Restore, "r", false, "Restore metric values from file before start")
	flag.StringVar(&conf.Key, "k", "", "Key to sign up data with SHA256 algorythm")
	flag.StringVar(&conf.ConnectionString, "d", "",
		"Database connection string(PostgreSql)") // postgresql://pguser:pgpwd@localhost:5432/devops?sslmode=disable
	flag.StringVar(&conf.CryptoKeyFile, "crypto-key", "", "Public Key file address")
	flag.Parse()

	err := env.Parse(&conf)
	if err != nil {
		return nil, fmt.Errorf("error reading agent  config: %w", err)
	}
	err = conf.parseConfigFie()
	if err != nil {
		return nil, fmt.Errorf("error reading agent config(Json): %w", err)
	}
	if err = conf.getCryptoKey(); err != nil {
		return nil, fmt.Errorf("error reading agent config, failed to get CryptoKey: %w", err)
	}

	return &conf, err
}
func (c *Config) UnmarshalJSON(data []byte) (err error) {
	var cfg struct {
		ServerAddress    string `json:"address"`
		StoreInterval    string `json:"store_interval"`
		StoreFile        string `json:"store_file"`
		Restore          bool   `json:"restore"`
		ConnectionString string `json:"database_dsn"`
		CryptoKeyFile    string `json:"crypto_key"`
	}

	if err = json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	c.ServerAddress = cfg.ServerAddress
	c.StoreInterval, err = time.ParseDuration(cfg.StoreInterval)
	if err != nil {
		return err
	}
	c.StoreFile = cfg.StoreFile
	c.Restore = cfg.Restore
	c.ConnectionString = cfg.ConnectionString
	c.CryptoKeyFile = cfg.CryptoKeyFile
	return err
}

func (c *Config) parseConfigFie() error {
	if c.confFile == "" {
		return nil
	}
	josnFile, err := os.Open(c.confFile)
	if err != nil {
		return err
	}
	defer josnFile.Close()

	r, err := io.ReadAll(josnFile)
	if err != nil {
		return err
	}
	var cfg Config
	err = json.Unmarshal(r, &cfg)
	if err != nil {
		return err
	}

	if c.ServerAddress == "" {
		c.ServerAddress = cfg.ServerAddress
	}
	if c.StoreInterval == time.Second*0 {
		c.StoreInterval = cfg.StoreInterval
	}
	if !c.Restore {
		c.Restore = cfg.Restore
	}
	if c.StoreFile == "" {
		c.StoreFile = cfg.StoreFile
	}
	if c.ConnectionString == "" {
		c.ConnectionString = cfg.ConnectionString
	}
	if c.CryptoKeyFile == "" {
		c.CryptoKeyFile = cfg.CryptoKeyFile
	}
	return nil
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
	_, err = file.Read(data)
	if err != nil {
		return err
	}
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
