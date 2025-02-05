package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

var Config struct {
	Addr        string `env:"SERVER_ADDRESS"`
	BaseURL     string `env:"BASE_URL"`
	FileStorage string `env:"FILE_STORAGE_PATH"`
	DataBase    string `env:"DATABASE_DSN"`
}

func InitConfig() error {
	flag.StringVar(&Config.Addr, "a", ":8080", "host and port to run server :port")
	flag.StringVar(&Config.BaseURL, "b", "http://localhost:8080", "static short url")
	flag.StringVar(&Config.FileStorage, "f", "tmp/short-url-db.json", "storage data")
	flag.StringVar(&Config.DataBase, "d", "", "name data base")
	flag.Parse()
	if err := env.Parse(&Config); err != nil {
		return err
	}
	return nil
}
