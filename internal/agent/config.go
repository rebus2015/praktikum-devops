// Package agent реализует агент сбора метрик
package agent

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

type Config struct {
	ServerAddress  string        `env:"ADDRESS" json:"address"`                 // Адрес сервера
	ReportInterval time.Duration `env:"PUSH_TIMEOUT" json:"report_interval"`    // Интервал отправки метрик на сервер
	PollInterval   time.Duration `env:"POLL_INTERVAL" json:"poll_interval"`     // Интервал сбора метрик
	Key            string        `env:"KEY" json:"-"`                           // Ключ для подписи данных
	RateLimit      int           `env:"RATE_LIMIT" json:"-"`                    // Количество одновременно исходящих запросов на сервер
	CryptoKeyFile  string        `env:"CRYPTO_KEY" json:"crypto_key,omitempty"` // Путь к файлу с открытым ключом
	confFile       string        `env:"CONFIG" json:"-"`
	CryptoKey      *rsa.PublicKey
}

func GetConfig() (*Config, error) {
	conf := Config{}
	flag.StringVar(&conf.confFile, "config", "", "Pass the conf.json path")
	flag.StringVar(&conf.confFile, "c", "", "Pass the conf.json path (shorthand)")
	flag.StringVar(&conf.ServerAddress, "a", "127.0.0.1:8080", "Server address")
	flag.DurationVar(&conf.ReportInterval, "r", time.Second*11, "Interval before push metrics to server")
	flag.DurationVar(&conf.PollInterval, "p", time.Second*5, "Interval between metrics reads from runtime")
	flag.StringVar(&conf.Key, "k", "", "Key to sign up data with SHA256 algorythm")
	flag.IntVar(&conf.RateLimit, "l", 5, "Workers count")
	flag.StringVar(&conf.CryptoKeyFile, "crypto-key", "", "Public Key file address") //"/Users/mak/go/praktikum-devops/keys/publicKey.pem",
	flag.Parse()

	err := env.Parse(&conf)
	if err != nil {
		return nil, fmt.Errorf("error reading agent  config(ENV): %w", err)
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
	var conf struct {
		ServerAddress  string `json:"address"`
		ReportInterval string `json:"report_interval"`
		PollInterval   string `json:"poll_interval"`
		CryptoKeyFile  string `json:"crypto_key"`
	}

	if err = json.Unmarshal(data, &conf); err != nil {
		return err
	}
	c.ServerAddress = conf.ServerAddress
	c.ReportInterval, err = time.ParseDuration(conf.ReportInterval)
	if err != nil {
		return err
	}
	c.PollInterval, err = time.ParseDuration(conf.PollInterval)
	if err != nil {
		return err
	}
	c.CryptoKeyFile = conf.CryptoKeyFile
	return err
}

func (c *Config) parseConfigFile() error {
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
	if c.ReportInterval == time.Second*0 {
		c.ReportInterval = cfg.ReportInterval
	}
	if c.PollInterval == time.Second*0 {
		c.PollInterval = cfg.PollInterval
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
	//1. Read the public key information and put it in the data variable
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

	//3. Use x509 to parse the encoded public key
	pubInterface, err2 := x509.ParsePKIXPublicKey(block.Bytes)
	if err2 != nil {
		return err2
	}
	c.CryptoKey = pubInterface.(*rsa.PublicKey)
	return nil
}
