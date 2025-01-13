package config

import (
	"flag"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Addr        string `env:"SERVER_ADDRESS"`
	BaseURL     string `env:"BASE_URL"`
	FileStorage string `env:"FILE_STORAGE_PATH"`
}

var Cfg *Config = &Config{}

func InitConfig() error {
	cfg := &Config{}
	addr := flag.String("a", ":8080", "host and port to run server :port")
	resURL := flag.String("b", "http://localhost:8080", "static short url")
	fStorage := flag.String("f", "", "storage data")
	flag.Parse()
	if err := env.Parse(cfg); err != nil {
		return err
	}

	if cfg.Addr == "" {
		cfg.Addr = *addr
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = *resURL
	}
	switch {
	case cfg.FileStorage != "":
		if err := validatePath(cfg); err != nil {
			return fmt.Errorf("validate path: %w", err)
		}
	case *fStorage != "":
		cfg.FileStorage = *fStorage
		if err := validatePath(cfg); err != nil {
			return fmt.Errorf("validate path: %w", err)
		}
	default:
		cfg.FileStorage = "without"
	}
	if err := validatePort(cfg.Addr); err != nil {
		return err
	}
	fmt.Println("----------", cfg.FileStorage)
	Cfg = cfg
	return nil
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

func validatePath(cfg *Config) error {
	pathTemp := cfg.FileStorage
	path := filepath.FromSlash(pathTemp)
	if filepath.Ext(path) == "" {
		path = fmt.Sprintf("%s.json", path)
	}
	cfg.FileStorage = path
	return nil
}
