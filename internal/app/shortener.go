package app

import (
	"fmt"
	"math/rand"
	"net/url"

	"github.com/hollgett/shortURL.git/internal/logger"
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

func (sh *Shortener) RandomID() string {
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
	if _, err := sh.Repo.Find(string(shortLink)); err != nil {
		return string(shortLink)
	} else {
		return sh.RandomID()
	}
}

// processing post request
func (sh *Shortener) CreateShortURL(requestData string) (string, error) {
	//checking link
	originalURL := requestData
	if err := isValidURL(originalURL); err != nil {
		logger.LogInfo("validated URL", zap.String("error", err.Error()))
		return "", fmt.Errorf("request URL doesn't match, error: %w", err)
	}
	//create short route
	shortLink := sh.RandomID()
	sh.Repo.Save(shortLink, originalURL)
	logger.LogInfo("data save storage", zap.String("key", shortLink), zap.String("value", originalURL))

	//return response
	return shortLink, nil
}

// processing post request
func (sh *Shortener) GetShortURL(pathURL string) (string, error) {
	//search exist short url and return original URL
	shortLink := pathURL

	originalURL, err := sh.Repo.Find(shortLink)
	logger.LogInfo("data find start", zap.String("key", shortLink))
	if err != nil {
		logger.LogInfo("data find", zap.String("error", err.Error()))
		return "", fmt.Errorf("find original link error: %w", err)
	}

	logger.LogInfo("data find complete", zap.String("value", originalURL))
	return originalURL, nil
}
