package api

import (
	"fmt"
	"net/http"

	"github.com/hollgett/shortURL.git/internal/app"
)

type HandlerAPI struct {
	ShortenerService *app.ShortenerHandler
}

func NewHandlerAPI(shortenerHandler *app.ShortenerHandler) *HandlerAPI {
	return &HandlerAPI{ShortenerService: shortenerHandler}
}

func (h *HandlerAPI) shortURLPost(w http.ResponseWriter, r *http.Request) {
	//create short url
	shortLink, err := h.ShortenerService.CreateShortURL(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//return response
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `http://localhost:8080/%s`, shortLink)
}

// processing post request
func (h *HandlerAPI) shortURLGet(w http.ResponseWriter, r *http.Request) {
	//search exist short url and return original URL
	originalURL, err := h.ShortenerService.GetShortURL(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (h *HandlerAPI) ShortURLmiddleware() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//checking http method and redirect
		switch r.Method {
		case http.MethodPost:
			h.shortURLPost(w, r)
		case http.MethodGet:
			h.shortURLGet(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}
