package api

import (
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/shortURL.git/internal/app"
	"github.com/hollgett/shortURL.git/internal/jsonutil"
	"github.com/hollgett/shortURL.git/internal/logger"
	"github.com/hollgett/shortURL.git/internal/models"
	"go.uber.org/zap"
)

type HandlerAPI struct {
	ShortenerService app.ShortenerHandler
}

func NewHandlerAPI(shortenerHandler app.ShortenerHandler) *HandlerAPI {
	return &HandlerAPI{ShortenerService: shortenerHandler}
}

func (h *HandlerAPI) HandlePlainTextRequest(w http.ResponseWriter, r *http.Request) {
	urlByte, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		ResponseWithError(w, "handlePlainTextRequest read body", err.Error(), http.StatusBadRequest)
		return
	}
	if len(urlByte) == 0 {
		ResponseWithError(w, "handlePlainTextRequest read body", "request body empty", http.StatusBadRequest)
		return
	}
	logger.LogInfo("body", zap.String("value", string(urlByte)))
	shLink, err := h.ShortenerService.CreateShortURL(strings.TrimSpace(string(urlByte)))
	if err != nil {
		ResponseWithError(w, "CreateShortURL", err.Error(), http.StatusBadRequest)
		return
	}
	ResponseWithSuccess(w, "Content-Type", "text/plain", shLink, http.StatusCreated)
}

func (h *HandlerAPI) HandleJSONRequest(w http.ResponseWriter, r *http.Request) {
	request := models.RequestJSON{}
	if err := jsonutil.DecodeJSON(r.Body, &request); err != nil {
		ResponseWithError(w, "handleJSONRequest json decode", err.Error(), http.StatusBadRequest)
		return
	}
	logger.LogInfo("handleJSONRequest json decode", zap.Any("argument", request))
	shLink, err := h.ShortenerService.CreateShortURL(strings.TrimSpace(request.RequestURL))
	if err != nil {
		ResponseWithError(w, "CreateShortURL", err.Error(), http.StatusBadRequest)
		return
	}
	ResponseWithSuccess(w, "Content-Type", "application/json", shLink, http.StatusCreated)
}

// processing post request
func (h *HandlerAPI) ShortURLGet(w http.ResponseWriter, r *http.Request) {
	//search exist short url and return original URL
	short := chi.URLParam(r, "short")
	logger.LogInfo("ShortURLGet take", zap.String("URL", short))

	originalURL, err := h.ShortenerService.GetShortURL(short)
	if err != nil {
		ResponseWithError(w, "ShortURLGet", err.Error(), http.StatusBadRequest)
		return
	}
	ResponseWithSuccess(w, "Location", originalURL, "", http.StatusTemporaryRedirect)
	logger.LogInfo("ShortURLGet complete", zap.String("location", originalURL))
}
