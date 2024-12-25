package api

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/shortURL.git/internal/app"
	"github.com/hollgett/shortURL.git/internal/config"
	"github.com/hollgett/shortURL.git/internal/logger"
	"github.com/hollgett/shortURL.git/internal/models"
	"go.uber.org/zap"
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
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		logger.Log.Info("error parse media type",
			zap.String("catch error", err.Error()),
		)
		http.Error(w, fmt.Sprint("request body error:", err.Error()), http.StatusBadRequest)
	}
	defer r.Body.Close()
	//create short url
	switch mediaType {
	case "text/plain", "application/x-www-form-urlencoded":
		logger.Log.Info(
			"handlePlainTextRequest start",
			zap.String("content type", mediaType),
		)
		//read request body plain text
		h.handlePlainTextRequest(w, r)
	case "application/json":
		logger.Log.Info(
			"handleJSONRequest start",
			zap.String("content type", mediaType),
		)
		//read request body json
		h.handleJSONRequest(w, r)
	default:
		logger.Log.Info("unsupported media type")
		http.Error(w, "unsupported media type", http.StatusBadRequest)
	}
}

func (h *HandlerAPI) handlePlainTextRequest(w http.ResponseWriter, r *http.Request) {
	urlByte, err := io.ReadAll(r.Body)
	originalURL := string(urlByte)
	if err != nil || len(urlByte) == 0 {
		logger.Log.Info(
			"handlePlainTextRequest read request body error",
			zap.String("request body", originalURL))
		http.Error(w, fmt.Sprintf("request body got: \"%s\"", originalURL), http.StatusBadRequest)
		return
	}

	logger.Log.Info("handlePlainTextRequest work",
		zap.String("transfer to CreateShortURL data", originalURL))
	shLink, err := h.ShortenerService.CreateShortURL(originalURL)
	if err != nil {
		logger.Log.Info(
			"handlePlainTextRequest work catch: CreateShortURL error",
			zap.String("catch error", err.Error()),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info(
		"handlePlainTextRequest work",
		zap.String("short URL got", shLink),
	)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	response := fmt.Sprintf("%s%s%s%s", "http://", r.Host, r.URL, shLink)
	fmt.Fprint(w, response)
	logger.Log.Info(
		"handlePlainTextRequest work complete",
		zap.String("response data", response),
	)
}

func (h *HandlerAPI) handleJSONRequest(w http.ResponseWriter, r *http.Request) {
	var request models.RequestJson
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.Log.Info(
			"handleJSONRequest work json decoder error",
			zap.String("error catch", err.Error()),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info(
		"handleJSONRequest transfer to CreateShortURL data",
		zap.Any("request decode", request),
	)
	shLink, err := h.ShortenerService.CreateShortURL(request.RequestURL)
	if err != nil {
		logger.Log.Info(
			"handleJSONRequest work catch: CreateShortURL error",
			zap.String("catch error", err.Error()),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := models.ResponseJson{
		ResponseURL: fmt.Sprintf("%s/%s", h.cfg.BaseURL, shLink),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&response)
	logger.Log.Info(
		"handleJSONRequest work complete",
		zap.Any("response data", response),
	)
}

// processing post request
func (h *HandlerAPI) ShortURLGet(w http.ResponseWriter, r *http.Request) {
	//search exist short url and return original URL
	short := chi.URLParam(r, "short")
	logger.Log.Info(
		"ShortURLGet take",
		zap.String("URL param", short),
	)
	originalURL, err := h.ShortenerService.GetShortURL(short)
	if err != nil {
		logger.Log.Info(
			"ShortURLGet catch error from GetShortURL",
			zap.String("error", err.Error()),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info(
		"ShortURLGet complete",
		zap.String("response header location set", originalURL),
	)
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
