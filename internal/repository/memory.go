package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hollgett/shortURL.git/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DataStorage struct {
	data map[string]string
}

func NewStorage() (Storage, error) {
	ds := DataStorage{
		data: make(map[string]string),
	}
	return &ds, nil
}

func (ds *DataStorage) Save(ctx context.Context, shortLink, originURL string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	ds.data[shortLink] = originURL
	return nil

}

func (ds *DataStorage) Find(ctx context.Context, shortLink string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if originURL, ok := ds.data[shortLink]; ok {
		return originURL, nil
	}
	return "", errors.New("the object does not exist in storage")

}

func (ds *DataStorage) Close() error {
	return nil
}

func (ds *DataStorage) Ping(context.Context) error {
	return nil
}

func (ds *DataStorage) SaveBatch(data []models.DBBatch) error {
	for _, v := range data {
		if err := ds.Save(context.TODO(), v.Short, v.Original); err != nil {
			return fmt.Errorf("save error: %w", err)
		}
	}
	return nil
}
