package app

import (
	"errors"
	"io"
	"math/rand"
	"net/url"

	"github.com/hollgett/shortURL.git/internal/repository"
)

type ShortenerHandler struct {
	Repo repository.Storage
}

func NewShortenerHandler(repo repository.Storage) *ShortenerHandler {
	return &ShortenerHandler{Repo: repo}
}

func isValid(u string) bool {
	_, err := url.Parse(u)
	return err == nil
}

func (sh *ShortenerHandler) randomID() string {
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
		return sh.randomID()
	}
}

// processing post request
func (sh *ShortenerHandler) CreateShortURL(body io.ReadCloser) (string, error) {
	//read request body
	defer body.Close()
	urlByte, err := io.ReadAll(body)
	if err != nil || len(urlByte) == 0 {
		return "", errors.Join(errors.New("request body error: "), err)
	}
	//checking link
	originalURL := string(urlByte)
	if !isValid(originalURL) {
		return "", errors.New("request link doesn't match")
	}
	//create short route
	shortLink := sh.randomID()
	sh.Repo.Save(shortLink, originalURL)
	//return response
	return shortLink, nil
}

// processing post request
func (sh *ShortenerHandler) GetShortURL(pathURL string) (string, error) {
	//search exist short url and return original URL
	shortLink := pathURL[1:]
	if len(shortLink) == 0 {
		return "", errors.New("URL path length zero")
	}
	originalURL, err := sh.Repo.Find(shortLink)
	if err != nil {
		return "", errors.Join(errors.New("find original link error: "), err)
	}
	return originalURL, nil
}
