package main

type config struct {
	ServerAddress string
}

func getConfig() *config {
	return &config{ServerAddress: "127.0.0.1:8080"}
}
