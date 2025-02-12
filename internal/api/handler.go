package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/shortURL.git/internal/app"
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
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
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
	shLink, err := h.ShortenerService.CreateShortURL(ctx, strings.TrimSpace(string(urlByte)))
	if err != nil {
		ResponseWithError(w, "CreateShortURL", err.Error(), http.StatusBadRequest)
		return
	}
	ResponseWithSuccess(w, "Content-Type", "text/plain", shLink, http.StatusCreated)
}

func (h *HandlerAPI) HandleJSONRequest(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	request := models.RequestJSON{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ResponseWithError(w, "handleJSONRequest json decode", err.Error(), http.StatusBadRequest)
		return
	}
	shLink, err := h.ShortenerService.CreateShortURL(ctx, strings.TrimSpace(request.RequestURL))
	if err != nil {
		ResponseWithError(w, "CreateShortURL", err.Error(), http.StatusBadRequest)
		return
	}
	ResponseWithSuccess(w, "Content-Type", "application/json", shLink, http.StatusCreated)
}

// processing post request
func (h *HandlerAPI) ShortURLGet(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	//search exist short url and return original URL
	short := chi.URLParam(r, "short")
	logger.Log.Info("--------------",zap.Any("dsada",short))
	originalURL, err := h.ShortenerService.GetShortURL(ctx, short)
	if err != nil {
		ResponseWithError(w, "ShortURLGet", err.Error(), http.StatusBadRequest)
		return
	}
	ResponseWithSuccess(w, "Location", originalURL, "", http.StatusTemporaryRedirect)
}

func (h *HandlerAPI) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := h.ShortenerService.Ping(ctx); err != nil {
		ResponseWithError(w, "db ping", err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *HandlerAPI) BatchReq(w http.ResponseWriter, r *http.Request) {
	var Data []models.RequestBatch
	if err := json.NewDecoder(r.Body).Decode(&Data); err != nil {
		ResponseWithError(w, "decoder body", err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("TEST", zap.Any("value", Data))

	respData, err := h.ShortenerService.ShortenBatch(Data)
	if err != nil {
		ResponseWithError(w, "shorten batch body", err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Log.Info("TEST", zap.Any("value", respData))
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(respData); err != nil {
		logger.Log.Info("encoder BatchReq", zap.Error(err))
	}
}
