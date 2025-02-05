package api

import (
	"fmt"
	"net/http"

	"github.com/hollgett/shortURL.git/internal/config"
	"github.com/hollgett/shortURL.git/internal/jsonutil"
	"github.com/hollgett/shortURL.git/internal/logger"
	"github.com/hollgett/shortURL.git/internal/models"
	"go.uber.org/zap"
)

func ResponseWithError(w http.ResponseWriter, logMess, err string, status int) {
	logger.LogInfo(logMess, zap.String("error", err))
	http.Error(w, err, status)
}

func ResponseWithSuccess(w http.ResponseWriter, headerK, headerV, shLink string, status int) {
	w.Header().Set(headerK, headerV)
	w.WriteHeader(status)
	if len(shLink) != 0 {
		switch headerV {
		case "application/json":
			response := models.ResponseJSON{
				ResponseURL: fmt.Sprintf("%s/%s", config.Cfg.BaseURL, shLink),
			}
			if err := jsonutil.EncodeJSON(w, response); err != nil {
				ResponseWithError(w, "json encode", err.Error(), http.StatusBadRequest)
				return
			}
		default:
			response := fmt.Sprintf("%s/%s", config.Cfg.BaseURL, shLink)
			fmt.Fprint(w, response)
		}
		logger.LogInfo("response server", zap.String("data", shLink))
	}
}
