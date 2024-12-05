package repository

import (
	"errors"
)

type Storage interface {
	Save(shortLink, originURL string)
	Find(shortLin string) (string, error)
}

type DataStorage struct {
	data map[string]string
}

func NewStorage() *DataStorage {
	return &DataStorage{
		data: make(map[string]string),
	}
}

func (ds *DataStorage) Save(shortLink, originURL string) {
	ds.data[shortLink] = originURL
}

func (ds *DataStorage) Find(shortLin string) (string, error) {
	if originURL, ok := ds.data[shortLin]; ok {
		return originURL, nil
	}
	return "", errors.New("the object does not exist in storage")
}