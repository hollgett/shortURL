package repository

import (
	"context"
	"errors"

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

func (ds *DataStorage) Save(shortLink, originURL string) error {
	ds.data[shortLink] = originURL
	return nil
}

func (ds *DataStorage) Find(shortLink string) (string, error) {
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
