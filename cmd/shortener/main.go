package main

import (
	"github.com/hollgett/shortURL.git/internal/api"
	"github.com/hollgett/shortURL.git/internal/app"
	"github.com/hollgett/shortURL.git/internal/config"
	"github.com/hollgett/shortURL.git/internal/logger"
	"github.com/hollgett/shortURL.git/internal/repository"
	"github.com/hollgett/shortURL.git/internal/server"
	"go.uber.org/zap"
)

func main() {
	err := logger.InitLogger()
	if err != nil {
		panic(err)
	}
	cfg := config.InitConfig()

	//data base
	repo := repository.NewStorage()
	//logic handler the short URl
	shortener := app.NewShortenerHandler(repo, cfg)
	//http logic
	apih := api.NewHandlerAPI(shortener, cfg)

	//start serve
	srv := server.NewServer(apih, cfg)
	if err := srv.ListenAndServe(); err != nil {
		logger.Log.Info(
			"listen and serve panic",
			zap.String("error", err.Error()),
		)
		panic(err)
	}
}
