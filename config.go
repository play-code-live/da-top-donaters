package main

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Host           string `env:"HOST"`
	Port           int    `env:"PORT"`
	SocketAddress  string `env:"SOCKET_ADDRESS"`
	DaClientId     string `env:"CLIENT_ID"`
	DaClientSecret string `env:"CLIENT_SECRET"`
}

func GetConfig() *Config {
	_ = godotenv.Load(".env.local", ".env")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	cfg := &Config{
		Host:           os.Getenv("HOST"),
		Port:           port,
		SocketAddress:  os.Getenv("SOCKET_ADDRESS"),
		DaClientId:     os.Getenv("CLIENT_ID"),
		DaClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	return cfg
}
