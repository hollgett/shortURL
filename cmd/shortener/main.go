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
	if err := logger.InitLogger(); err != nil {
		panic(err)
	}
	if err := config.InitConfig(); err != nil {
		logger.LogInfo("init config", zap.String("error", err.Error()))
		return
	}
	logger.LogInfo("server start", zap.Any("arg", config.Cfg))
	//data base
	repo := repository.NewStorage()
	//logic handler the short URl
	shortener := app.NewShortenerHandler(repo)
	//http logic
	apih := api.NewHandlerAPI(shortener)

	//start serve
	srv := server.NewServer(apih)
	if err := srv.ListenAndServe(); err != nil {
		logger.LogInfo("serve start", zap.String("error", err.Error()))
		panic(err)
	}
}
