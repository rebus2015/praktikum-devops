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
	"log"
	"os"
	"time"

	"github.com/caarlos0/env"
)

// Config хранит получныые занчеия конфигурации.
type Config struct {
	CryptoKey        *rsa.PrivateKey
	ServerAddress    string        `env:"ADDRESS" json:"address"`
	StoreFile        string        `env:"STORE_FILE" json:"store_file"`     // пустое значние отключает запись на диск
	Key              string        `env:"KEY"`                              // Ключ для создания подписи сообщения
	ConnectionString string        `env:"DATABASE_DSN" json:"database_dsn"` // Cтрока подключения к БД
	CryptoKeyFile    string        `env:"CRYPTO_KEY" json:"crypto_key"`     // путь к файлу с приватным ключом
	confFile         string        `env:"CONFIG" json:"-"`
	StoreInterval    time.Duration `env:"STORE_INTERVAL" json:"store_interval"` // 0 - синхронная запись
	Restore          bool          `env:"RESTORE" json:"restore"`               // загружать начальные значениея из файла

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
	err = conf.parseConfigFile()
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
		ConnectionString string `json:"database_dsn"`
		CryptoKeyFile    string `json:"crypto_key"`
		Restore          bool   `json:"restore"`
	}

	if err = json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("json.unmarshal error: %w", err)
	}
	c.ServerAddress = cfg.ServerAddress
	c.StoreInterval, err = time.ParseDuration(cfg.StoreInterval)
	if err != nil {
		return fmt.Errorf("time.ParseDuration error: %w", err)
	}
	c.StoreFile = cfg.StoreFile
	c.Restore = cfg.Restore
	c.ConnectionString = cfg.ConnectionString
	c.CryptoKeyFile = cfg.CryptoKeyFile
	return nil
}

func (c *Config) parseConfigFile() error {
	if c.confFile == "" {
		return nil
	}
	josnFile, err := os.Open(c.confFile)
	if err != nil {
		return fmt.Errorf("os.Open error: %w", err)
	}
	defer func() {
		err := josnFile.Close()
		if err != nil {
			log.Fatalf("error josnFile.Close : %v", err)
		}
	}()

	r, err := io.ReadAll(josnFile)
	if err != nil {
		return fmt.Errorf("io.ReadAll(josnFile) error: %w", err)
	}
	var cfg Config
	err = json.Unmarshal(r, &cfg)
	if err != nil {
		return fmt.Errorf("json.Unmarshal config error: %w", err)
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
	// 1. Read the private key information and put it in the data variable
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("os.Open(%s) error: %w", filename, err)
	}
	stat, err := file.Stat() // get file attribute information
	if err != nil {
		return fmt.Errorf("read file '%s' attributes error: %w", filename, err)
	}
	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		return fmt.Errorf("os.Read error: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("file close error: %w", err)
	}
	// 2. Decode the resulting string pem
	block, _ := pem.Decode(data)
	if block == nil {
		return errors.New("error reading key bytes")
	}

	privateKey, err3 := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err3 != nil {
		return fmt.Errorf("error reading agent  config: %w", err3)
	}
	c.CryptoKey = privateKey
	return nil
}
