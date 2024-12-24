package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/shortURL.git/internal/api"
	"github.com/hollgett/shortURL.git/internal/config"
)

// start serve
func NewServer(handler *api.HandlerAPI, config *config.Config) *http.Server {
	rtr := chi.NewMux()
	fmt.Println(config.Addr, config.BaseURL)
	rtr.Post("/", handler.ShortURLPost)
	rtr.Get("/{short}", handler.ShortURLGet)
	rtr.Post("/api/shorten", handler.ShortURLPost)
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Addr),
		Handler: rtr,
	}
}
