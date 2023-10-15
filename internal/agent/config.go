// Package agent реализует агент сбора метрик.
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
	log "github.com/sirupsen/logrus"
)

var (
	defReportInterval = time.Second * 5
	defPollInterval   = time.Second * 3
)

type Config struct {
	CryptoKey        *rsa.PublicKey
	ServerAddress    string        `env:"ADDRESS" json:"address"`                 // Адрес сервера
	CryptoKeyFile    string        `env:"CRYPTO_KEY" json:"crypto_key,omitempty"` // Путь к файлу с открытым ключом
	confFile         string        `env:"CONFIG" json:"-"`
	Key              string        `env:"KEY" json:"-"` // Ключ для подписи данных
	RPCServerAddress string        `env:"RPC_HOST" json:"-"`
	ReportInterval   time.Duration `env:"PUSH_TIMEOUT" json:"report_interval"` // Интервал отправки метрик на сервер
	PollInterval     time.Duration `env:"POLL_INTERVAL" json:"poll_interval"`  // Интервал сбора метрик
	RateLimit        int           `env:"RATE_LIMIT" json:"-"`                 // Количество одновременных запросов
	UseRPC           bool          `env:"RPC" json:"-"`                        // запускать gRPC- клиент

}

func GetConfig() (*Config, error) {
	conf := Config{}
	flag.StringVar(&conf.confFile, "config", "", "Pass the conf.json path")
	flag.StringVar(&conf.confFile, "c", "", "Pass the conf.json path (shorthand)")
	flag.StringVar(&conf.ServerAddress, "a", "127.0.0.1:8080", "Server address")
	flag.StringVar(&conf.RPCServerAddress, "tcp-host", ":3200", "Server host addr for RPC")
	flag.DurationVar(&conf.ReportInterval, "r", defReportInterval, "Interval before push metrics to server")
	flag.DurationVar(&conf.PollInterval, "p", defPollInterval, "Interval between metrics reads from runtime")
	flag.StringVar(&conf.Key, "k", "", "Key to sign up data with SHA256 algorythm")
	flag.IntVar(&conf.RateLimit, "l", 5, "Workers count")
	flag.StringVar(&conf.CryptoKeyFile, "crypto-key", "", "Public Key file address")
	flag.BoolVar(&conf.UseRPC, "grpc", true, "Start gRPC client for metrics Update")
	flag.Parse()

	err := env.Parse(&conf)
	if err != nil {
		return nil, fmt.Errorf("error reading agent  config(ENV): %w", err)
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
	var conf struct {
		ServerAddress  string `json:"address"`
		ReportInterval string `json:"report_interval"`
		PollInterval   string `json:"poll_interval"`
		CryptoKeyFile  string `json:"crypto_key"`
	}

	if err = json.Unmarshal(data, &conf); err != nil {
		return fmt.Errorf("unmarshal config failed with error: ,%w", err)
	}
	c.ServerAddress = conf.ServerAddress
	c.ReportInterval, err = time.ParseDuration(conf.ReportInterval)
	if err != nil {
		return fmt.Errorf("time.ParceDuration for ReportInterval failed with error: ,%w", err)
	}
	c.PollInterval, err = time.ParseDuration(conf.PollInterval)
	if err != nil {
		return fmt.Errorf("time.ParceDuration PollInterval failed with error: ,%w", err)
	}
	c.CryptoKeyFile = conf.CryptoKeyFile
	return nil
}

func (c *Config) parseConfigFile() error {
	if c.confFile == "" {
		return nil
	}
	josnFile, err := os.Open(c.confFile)
	if err != nil {
		return fmt.Errorf("open json config file failed with error: ,%w", err)
	}
	defer func() {
		err := josnFile.Close()
		if err != nil {
			log.Printf("failed to close json config file: %v", err.Error())
		}
	}()

	r, err := io.ReadAll(josnFile)
	if err != nil {
		return fmt.Errorf("failed to read json config file: %v", err.Error())
	}
	var cfg Config
	err = json.Unmarshal(r, &cfg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json config file: %v", err.Error())
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
	// 1. Read the public key information and put it in the data variable
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error trying to open json config file: %w", err)
	}
	stat, err := file.Stat() // Get file attribute information
	if err != nil {
		return fmt.Errorf("error trying get file attribute information for json config file: %w", err)
	}
	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		return fmt.Errorf("error trying to read json config file: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("error trying to close json config file: %w", err)
	}
	// 2. Decode the resulting string pem
	block, _ := pem.Decode(data)

	// 3. Use x509 to parse the encoded public key
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("error trying use x509 to parse the encoded public key: %w", err)
	}
	cKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("error for type assertion pubInterface.(*rsa.PublicKey)")
	}
	c.CryptoKey = cKey
	return nil
}
