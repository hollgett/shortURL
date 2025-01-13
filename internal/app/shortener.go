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
	logger.LogInfo("CreateShortURL start", zap.String("value", requestData))
	//checking link
	if err := isValidURL(requestData); err != nil {
		return "", fmt.Errorf("request URL doesn't match, error: %w", err)
	}
	//create short route
	shortLink := sh.RandomID()
	sh.Repo.Save(shortLink, requestData)
	//return response
	logger.LogInfo("CreateShortURL complete", zap.String("result", shortLink))
	return shortLink, nil
}

// processing post request
func (sh *Shortener) GetShortURL(pathURL string) (string, error) {
	logger.LogInfo("GetShortURL start", zap.String("value", pathURL))
	//search exist short url and return original URL
	originalURL, err := sh.Repo.Find(pathURL)
	if err != nil {
		logger.LogInfo("data find", zap.String("error", err.Error()))
		return "", fmt.Errorf("find original link error: %w", err)
	}
	logger.LogInfo("GetShortURL complete", zap.String("result", originalURL))
	return originalURL, nil
}
