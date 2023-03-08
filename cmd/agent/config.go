package main

import (
	"time"
)

type config struct {
	ServerAddress  string
	ReportInternal time.Duration
	PollInterval   time.Duration
	GaugeList      []string
}

func getConfig() *config {
	return &config{
		ServerAddress:  "127.0.0.1:8080",
		ReportInternal: 2 * time.Second,
		PollInterval:   5 * time.Second,
	}
}
