package app

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"

	"github.com/hollgett/shortURL.git/internal/logger"
	"github.com/hollgett/shortURL.git/internal/models"
	"github.com/hollgett/shortURL.git/internal/repository"
	"go.uber.org/zap"
)

type Shortener struct {
	Repo repository.Storage
}

func NewShortenerHandler(repo repository.Storage) ShortenerHandler {
	return &Shortener{Repo: repo}
}

func isValidURL(URL string) error {
	_, err := url.Parse(URL)
	return err
}

func (sh *Shortener) RandomID(ctx context.Context) string {
	//create random URL path
	shortLink := make([]byte, 8)
	for i := range shortLink {
		if rand.Intn(2) == 0 {
			shortLink[i] = byte(rand.Intn(26) + 65)
		} else {
			shortLink[i] = byte(rand.Intn(26) + 97)
		}
	}
	// search exist link
	if _, err := sh.Repo.Find(ctx, string(shortLink)); err != nil {
		return string(shortLink)
	} else {
		return sh.RandomID(ctx)
	}
}

// processing post request
func (sh *Shortener) CreateShortURL(ctx context.Context, requestData string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	logger.LogInfo("CreateShortURL start", zap.String("value", requestData))
	//checking link
	if err := isValidURL(requestData); err != nil {
		return "", fmt.Errorf("request URL doesn't match, error: %w", err)
	}
	//create short route
	shortLink := sh.RandomID(ctx)
	if err := sh.Repo.Save(ctx, shortLink, requestData); err != nil {
		return "", fmt.Errorf("save data, error: %w", err)
	}

	//return response
	logger.LogInfo("CreateShortURL complete", zap.String("result", shortLink))
	return shortLink, nil
}

// processing post request
func (sh *Shortener) GetShortURL(ctx context.Context, pathURL string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	logger.LogInfo("GetShortURL start", zap.String("value", pathURL))
	//search exist short url and return original URL
	originalURL, err := sh.Repo.Find(ctx, pathURL)
	if err != nil {
		logger.LogInfo("data find", zap.String("error", err.Error()))
		return "", fmt.Errorf("find original link error: %w", err)
	}
	logger.LogInfo("GetShortURL complete", zap.String("result", originalURL))
	return originalURL, nil
}

func (sh *Shortener) Ping(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return sh.Repo.Ping(ctx)
}

func (sh *Shortener) ShortenBatch(original []models.RequestBatch) ([]models.ResponseBatch, error) {
	var dbData []models.DBBatch
	for _, v := range original {
		dbData = append(dbData, models.DBBatch{
			Short:    sh.RandomID(context.TODO()),
			CorrId:   v.CorrId,
			Original: v.Original,
		})
	}
	if err := sh.Repo.SaveBatch(dbData); err != nil {
		return nil, fmt.Errorf("save batch: %w", err)
	}
	var respData []models.ResponseBatch
	for _, v := range dbData {
		respData = append(respData, models.ResponseBatch{
			CorrId: v.CorrId,
			Short:  v.Short,
		})
	}
	return respData, nil
}
