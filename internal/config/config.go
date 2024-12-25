package config

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/hollgett/shortURL.git/internal/logger"
	"go.uber.org/zap"
)

type Config struct {
	Addr    string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
}

func validatePort(addr string) error {
	listHP := strings.Split(addr, ":")
	if len(listHP) != 2 || listHP[1] == "" {
		return fmt.Errorf("address must be in the format :port, got: %s", addr)
	}

	if _, err := strconv.Atoi(listHP[1]); err != nil {
		return fmt.Errorf("invalid port number: %s", listHP[1])
	}

	return nil
}

func InitConfig() *Config {
	cfg := &Config{}

	addr := flag.String("a", ":8080", "host and port to run server :port")
	resURL := flag.String("b", "http://localhost:8080", "static short url")
	flag.Parse()
	if err := env.Parse(cfg); err != nil {
		logger.Log.Info("env parse panic",
			zap.String("env parse error", err.Error()),
		)
		panic(err)
	}
	if cfg.Addr == "" || cfg.BaseURL == "" {
		if cfg.Addr == "" {
			cfg.Addr = *addr
		}
		if cfg.BaseURL == "" {
			cfg.BaseURL = *resURL
		}
	}
	if err := validatePort(cfg.Addr); err != nil {
		logger.Log.Info("config panic",
			zap.String("config address", cfg.Addr),
			zap.String("validate port panic", err.Error()),
		)
		panic(err)
	}
	logger.Log.Info(
		"server start with config",
		zap.String("server address", cfg.Addr),
		zap.String("server base URL", cfg.BaseURL),
	)
	return cfg
}
