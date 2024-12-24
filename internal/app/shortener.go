package app

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/url"

	"github.com/hollgett/shortURL.git/internal/config"
	"github.com/hollgett/shortURL.git/internal/repository"
)

type Data struct {
	Url string `json:"url"`
}

//go:generate mockgen -source=./shortener.go -destination=../mock/shortener.go -package=mock
type ShortenerHandler interface {
	RandomID() string
	CreateShortURL(body string) (string, error)
	CreateShortURLjson(body io.ReadCloser) (string, error)
	GetShortURL(pathURL string) (string, error)
}

type Shortener struct {
	Repo   repository.Storage
	config *config.Config
}

func NewShortenerHandler(repo repository.Storage, config *config.Config) ShortenerHandler {
	return &Shortener{Repo: repo, config: config}
}

func isValid(u string) bool {
	_, err := url.Parse(u)
	return err == nil
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
func (sh *Shortener) CreateShortURL(body string) (string, error) {
	//checking link
	originalURL := body
	if !isValid(originalURL) {
		return "", errors.New("request link doesn't match")
	}
	//create short route
	shortLink := sh.RandomID()
	sh.Repo.Save(shortLink, originalURL)
	//return response
	return shortLink, nil
}

func (sh *Shortener) CreateShortURLjson(body io.ReadCloser) (string, error) {
	var d Data
	//read request body
	err := json.NewDecoder(body).Decode(&d)
	if err != nil {
		return "", errors.Join(errors.New("request body error: "), err)
	}
	//checking link
	originalURL := d.Url
	if !isValid(originalURL) {
		return "", errors.New("request link doesn't match")
	}
	//create short route
	sh.Repo.Save("test", originalURL)
	//return response
	return "", nil
}

// processing post request
func (sh *Shortener) GetShortURL(pathURL string) (string, error) {
	//search exist short url and return original URL
	shortLink := pathURL
	originalURL, err := sh.Repo.Find(shortLink)
	if err != nil {
		return "", errors.Join(errors.New("find original link error: "), err)
	}
	return originalURL, nil
}
