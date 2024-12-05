package server

import (
	"net/http"

	"github.com/hollgett/shortURL.git/internal/api"
)

// start serve
func NewServer(handler *api.HandlerAPI) *http.Server {
	rtr := http.NewServeMux()
	rtr.HandleFunc("/", handler.ShortURLmiddleware())

	return &http.Server{
		Addr:    `:8080`,
		Handler: rtr,
	}
}
