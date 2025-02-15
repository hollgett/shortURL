package repository

import (
	"context"

	"github.com/hollgett/shortURL.git/internal/models"
)

type Storage interface {
	Save(ctx context.Context, shortLink, originURL string) error
	Find(ctx context.Context, shortLink string) (string, error)
	Close() error
	Ping(context.Context) error
	SaveBatch(data []models.DBBatch) error
}
