package app

import (
	"fmt"
	"math/rand"
	"net/url"

	"github.com/hollgett/shortURL.git/internal/config"
	"github.com/hollgett/shortURL.git/internal/logger"
	"github.com/hollgett/shortURL.git/internal/repository"
	"go.uber.org/zap"
)

type Shortener struct {
	Repo   repository.Storage
	config *config.Config
}

func NewShortenerHandler(repo repository.Storage, config *config.Config) ShortenerHandler {
	return &Shortener{Repo: repo, config: config}
}

func isValidURL(URL string) error {
	logger.Log.Info("validated URL start",
		zap.String("url take", URL),
	)
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
		logger.Log.Info("validated URL complete with error",
			zap.String("error", err.Error()),
		)
		return "", fmt.Errorf("request URL doesn't match, error: %w", err)
	}
	//create short route
	shortLink := sh.RandomID()
	sh.Repo.Save(shortLink, originalURL)
	logger.Log.Info("data save storage",
		zap.String("data key save", shortLink),
		zap.String("data value save", originalURL),
	)
	//return response
	return shortLink, nil
}

// processing post request
func (sh *Shortener) GetShortURL(pathURL string) (string, error) {
	//search exist short url and return original URL
	shortLink := pathURL
	originalURL, err := sh.Repo.Find(shortLink)
	logger.Log.Info("data find storage start",
		zap.String("data key", shortLink),
	)
	if err != nil {
		logger.Log.Info("data find storage complete with error",
			zap.String("error", err.Error()),
		)
		return "", fmt.Errorf("find original link error: %w", err)
	}
	logger.Log.Info("data find storage complete",
		zap.String("data value find", originalURL),
	)
	return originalURL, nil
}
