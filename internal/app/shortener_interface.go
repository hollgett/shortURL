package app

import (
	"context"

	"github.com/hollgett/shortURL.git/internal/models"
)

//go:generate mockgen -source=./shortener_interface.go -destination=../mocks/shortener.go -package=mocks
type ShortenerHandler interface {
	RandomID(ctx context.Context) string
	CreateShortURL(ctx context.Context, requestData string) (string, error)
	GetShortURL(ctx context.Context, pathURL string) (string, error)
	ShortenBatch(original []models.RequestBatch) ([]models.ResponseBatch, error)
	Ping(ctx context.Context) error
}
