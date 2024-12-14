package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/shortURL.git/internal/app"
	"github.com/hollgett/shortURL.git/internal/config"
)

type HandlerAPI struct {
	ShortenerService app.ShortenerHandler
	cfg              *config.Config
}

func NewHandlerAPI(shortenerHandler app.ShortenerHandler, cfg *config.Config) *HandlerAPI {
	return &HandlerAPI{ShortenerService: shortenerHandler,
		cfg: cfg}
}

func (h *HandlerAPI) ShortURLPost(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Content-Type")
	//create short url
	switch header {
	case "text/plain; charset=utf-8", "text/plain":
		//read request body
		urlByte, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil || len(urlByte) == 0 {
			http.Error(w, fmt.Sprint("request body error:", err.Error()), http.StatusBadRequest)
			return
		}
		shLink, err := h.ShortenerService.CreateShortURL(string(urlByte))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "%s%s%s%s", "http://", r.Host, r.URL, shLink)
	case "application/json":
		var request RequestJson
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		shLink, err := h.ShortenerService.CreateShortURL(request.RequestURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response := ResponseJson{
			ResponseURL: fmt.Sprintf("%s/%s", h.cfg.BaseURL, shLink),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&response)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

// processing post request
func (h *HandlerAPI) ShortURLGet(w http.ResponseWriter, r *http.Request) {
	//search exist short url and return original URL
	short := chi.URLParam(r, "short")
	originalURL, err := h.ShortenerService.GetShortURL(short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
