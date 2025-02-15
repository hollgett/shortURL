package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/hollgett/shortURL.git/internal/api"
	"github.com/hollgett/shortURL.git/internal/app"
	"github.com/hollgett/shortURL.git/internal/config"
	"github.com/hollgett/shortURL.git/internal/logger"
	"github.com/hollgett/shortURL.git/internal/repository"
	"github.com/hollgett/shortURL.git/internal/server"
	"go.uber.org/zap"
)

func main() {
	//main context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handlerShutDown(cancel)

	//initialisation logger function
	if err := logger.InitLogger(); err != nil {
		logger.LogInfo("init logger error", zap.Error(err))
		cancel()
	}

	//initialisation config, read env and flag
	if err := config.InitConfig(); err != nil {
		logger.LogInfo("init config error", zap.Error(err))
		cancel()
	}
	logger.LogInfo("server config", zap.Any("arg", config.Config))

	//init repository
	repo := initStorage(ctx)
	defer func() {
		if err := repo.Close(); err != nil {
			logger.LogInfo("repo close", zap.Error(err))
		}
	}()

	//logic handler the shortener
	shortener := app.NewShortenerHandler(repo)

	//api logic
	api := api.NewHandlerAPI(shortener)

	//start serve with config
	srv := server.NewServer(api)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.LogInfo("server", zap.Error(err))
			cancel()
		}
	}()

	defer func() {
		if err := srv.Close(); err != nil {
			logger.LogInfo("server close", zap.Error(err))
		}
	}()

	<-ctx.Done()
}

func handlerShutDown(cancel context.CancelFunc) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	sig := <-signalCh
	logger.LogInfo("close app", zap.Any("signal", sig))
	cancel()

}

func initStorage(ctx context.Context) repository.Storage {
	//repository create
	var repo repository.Storage
	switch {
	case config.Config.DataBase != "":
		db, err := repository.NewDataBase(ctx,config.Config.DataBase)
		if err != nil {
			logger.LogInfo("storage start error", zap.Error(err))
			panic(err)
		}
		logger.LogInfo("postgres DB")
		repo = db
	case config.Config.FileStorage != "":
		fs, err := repository.NewFileStorage(config.Config.FileStorage)
		if err != nil {
			logger.LogInfo("file storage error", zap.Error(err))
			panic(err)
		}
		logger.LogInfo("file DB")
		repo = fs
	default:
		mem, err := repository.NewStorage()
		if err != nil {
			logger.LogInfo("storage start error", zap.Error(err))
			panic(err)
		}
		logger.LogInfo("memory DB")
		repo = mem
	}
	return repo
}
