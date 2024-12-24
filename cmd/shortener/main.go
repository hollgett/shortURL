package main

import (
	"github.com/hollgett/shortURL.git/internal/api"
	"github.com/hollgett/shortURL.git/internal/app"
	"github.com/hollgett/shortURL.git/internal/config"
	"github.com/hollgett/shortURL.git/internal/repository"
	"github.com/hollgett/shortURL.git/internal/server"
)

func main() {
	// cfg := config.InitConfig()
	cfg := &config.Config{
		Addr:    "8080",
		BaseURL: "http://localhost:8080",
	}
	//data base
	repo := repository.NewStorage()
	//logic handler the short URl
	shortener := app.NewShortenerHandler(repo, cfg)
	//http logic
	apih := api.NewHandlerAPI(shortener, cfg)

	//start serve
	srv := server.NewServer(apih, cfg)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
