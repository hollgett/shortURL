package config

import (
	"flag"
	"strconv"
	"strings"
)

type Config struct {
	Addr    string
	BaseURL string
}

func InitConfig() *Config {
	addr := flag.String("a", ":8080", "host and port to run server :port")
	resURL := flag.String("b", "http://localhost:8080", "static short url")

	port := *addr
	listHP := strings.Split(port, ":")
	if _, err := strconv.Atoi(listHP[1]); err != nil {
		panic("receiver port error: " + port)
	}

	// flag.Parse()
	return &Config{
		Addr:    listHP[1],
		BaseURL: *resURL,
	}
}
