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
	//initialisation logger function
	if err := logger.InitLogger(); err != nil {
		panic(err)
	}
	//initialisation config, read env and flag
	if err := config.InitConfig(); err != nil {
		logger.LogInfo("init config error", zap.Error(err))
		return
	}
	logger.LogInfo("server start", zap.Any("arg", config.Cfg))
	//data base create
	repo, err := repository.NewStorage()
	if err != nil {
		logger.LogInfo("storage start error", zap.Error(err))
		return
	}
	//logic handler the shortener
	shortener := app.NewShortenerHandler(repo)
	//api logic
	apih := api.NewHandlerAPI(shortener)

	//start serve with config
	srv := server.NewServer(apih)
	if err := srv.ListenAndServe(); err != nil {
		logger.LogInfo("serve start", zap.String("error", err.Error()))
		panic(err)
	}
}
