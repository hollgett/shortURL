package repository

import (
	"errors"

	"github.com/hollgett/shortURL.git/internal/config"
	"github.com/hollgett/shortURL.git/internal/logger"
	"go.uber.org/zap"
)

type DataStorage struct {
	data        map[string]string
	fileStorage bool
}

func NewStorage() (*DataStorage, error) {
	base := make(map[string]string)
	switch {
	case config.Cfg.FileStorage == "without":

		return &DataStorage{
			data:        base,
			fileStorage: false,
		}, nil
	default:
		if err := readFileStorage(&base); err != nil {
			return nil, err
		}
		logger.LogInfo("file storage data", zap.Int("count", len(base)))
		return &DataStorage{
			data:        base,
			fileStorage: true,
		}, nil
	}

}

func (ds *DataStorage) Save(shortLink, originURL string) {
	if ds.fileStorage {
		writeFileStorage(shortLink, originURL)
	}
	ds.data[shortLink] = originURL
}

func (ds *DataStorage) Find(shortLink string) (string, error) {
	if originURL, ok := ds.data[shortLink]; ok {
		return originURL, nil
	}
	return "", errors.New("the object does not exist in storage")
}
