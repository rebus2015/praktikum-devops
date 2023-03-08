package main

import (
	"time"
)

type config struct {
	ServerAddress  string
	ReportInternal time.Duration
	PollInterval   time.Duration
}

func getConfig() *config {
	return &config{
		ServerAddress:  "127.0.0.1:8080",
		PollInterval:   2 * time.Second,
		ReportInternal: 5 * time.Second,
	}
}
