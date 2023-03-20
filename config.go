package main

import (
	"github.com/golobby/dotenv"
	"os"
)

type Config struct {
	Host           string `env:"HOST"`
	Port           int    `env:"PORT"`
	SocketAddress  string `env:"SOCKET_ADDRESS"`
	DaClientId     string `env:"CLIENT_ID"`
	DaClientSecret string `env:"CLIENT_SECRET"`
}

func GetConfig() (*Config, error) {
	cfg := &Config{}

	file, err := os.Open(".env")
	if err != nil {
		return nil, err
	}

	return cfg, dotenv.NewDecoder(file).Decode(cfg)
}
