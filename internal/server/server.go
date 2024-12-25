package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/shortURL.git/internal/api"
	"github.com/hollgett/shortURL.git/internal/config"
	"github.com/hollgett/shortURL.git/internal/logger"
)

func setupRouters(handler *api.HandlerAPI) *chi.Mux {
	r := chi.NewMux()

	r.Use(logger.RequestMiddleware)
	r.Use(logger.ResponseMiddleware)
	r.Post("/", handler.ShortURLPost)
	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", handler.ShortURLPost)
	})
	r.Get("/{short}", handler.ShortURLGet)

	return r
}

// start serve
func NewServer(handler *api.HandlerAPI, config *config.Config) *http.Server {
	rtr := setupRouters(handler)
	return &http.Server{
		Addr:    config.Addr,
		Handler: rtr,
	}
}
