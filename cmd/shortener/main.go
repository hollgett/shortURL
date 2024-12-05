package main

import (
	"github.com/hollgett/shortURL.git/internal/api"
	"github.com/hollgett/shortURL.git/internal/app"
	"github.com/hollgett/shortURL.git/internal/repository"
	"github.com/hollgett/shortURL.git/internal/server"
)

func main() {
	//data base
	repo := repository.NewStorage()
	//logic handler the short URl
	shortener := app.NewShortenerHandler(repo)
	//http logic
	apih := api.NewHandlerAPI(shortener)
	//start serve
	srv := server.NewServer(apih)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
