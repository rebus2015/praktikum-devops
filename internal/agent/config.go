// Package agent реализует агент сбора метрик
package agent

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

type Config struct {
	ServerAddress  string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"PUSH_TIMEOUT"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
	RateLimit      int           `env:"RATE_LIMIT"` // Количество одновременно исходящих запросов на сервер
	CryptoKeyFile  string        `env:"CRYPTO_KEY"` // путь к файлу с открытым ключом
	CryptoKey      *rsa.PublicKey
}

func GetConfig() (*Config, error) {
	conf := Config{}
	flag.StringVar(&conf.ServerAddress, "a", "127.0.0.1:8080", "Server address")
	flag.DurationVar(&conf.ReportInterval, "r", time.Second*11, "Interval before push metrics to server")
	flag.DurationVar(&conf.PollInterval, "p", time.Second*5, "Interval between metrics reads from runtime")
	flag.StringVar(&conf.Key, "k", "", "Key to sign up data with SHA256 algorythm")
	flag.IntVar(&conf.RateLimit, "l", 5, "Workers count")
	flag.StringVar(&conf.CryptoKeyFile, "crypto-key", "/Users/mak/go/praktikum-devops/keys/publicKey.pem", "Public Key file address")
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
	//1. Read the public key information and put it in the data variable
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

	//3. Use x509 to parse the encoded public key
	pubInterface, err2 := x509.ParsePKIXPublicKey(block.Bytes)
	if err2 != nil {
		return err2
	}
	c.CryptoKey = pubInterface.(*rsa.PublicKey)
	return nil
}
