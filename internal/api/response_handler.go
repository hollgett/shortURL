package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hollgett/shortURL.git/internal/config"
	"github.com/hollgett/shortURL.git/internal/logger"
	"github.com/hollgett/shortURL.git/internal/models"
	"go.uber.org/zap"
)

func ResponseWithError(w http.ResponseWriter, logMess, err string, status int) {
	logger.LogInfo(logMess, zap.String("error", err))
	http.Error(w, err, status)
}

func ResponseWithSuccessJSON(w http.ResponseWriter, shLink string, status int) {
	response := models.ResponseJSON{
		ResponseURL: fmt.Sprintf("%s/%s", config.Config.BaseURL, shLink),
	}
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(response); err != nil {
		ResponseWithError(w, "json encode", err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b.Bytes())
	logger.LogInfo("response server", zap.String("data", shLink))
}

func ResponseWithSuccessText(w http.ResponseWriter, shLink string, status int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "%s/%s", config.Config.BaseURL, shLink); err != nil {
		ResponseWithError(w, "json encode", err.Error(), http.StatusInternalServerError)
		return
	}
	logger.LogInfo("response server", zap.String("data", shLink))
}

func ResponseWithSuccessGet(w http.ResponseWriter, original string) {
	w.Header().Set("Location", original)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func ResponseWithSuccessBatch(w http.ResponseWriter, respData []models.ResponseBatch) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(respData); err != nil {
		ResponseWithError(w, "json encode batch", err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
	logger.LogInfo("response batch success")
}
