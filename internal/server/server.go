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
	r.Use(api.CompressMiddleware)
	r.Post("/", handler.HandlePlainTextRequest)
	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", handler.HandleJSONRequest)
	})
	r.Get("/{short}", handler.ShortURLGet)

	return r
}

// start serve
func NewServer(handler *api.HandlerAPI) *http.Server {
	rtr := setupRouters(handler)
	return &http.Server{
		Addr:    config.Cfg.Addr,
		Handler: rtr,
	}
}
