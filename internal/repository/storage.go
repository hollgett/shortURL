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

func NewStorage() (Storage, error) {
	ds := DataStorage{
		data: make(map[string]string),
	}
	switch {
	case config.Cfg.FileStorage == "":
		ds.data = make(map[string]string)
		return &ds, nil
	default:
		fStorage, err := newFileStorage(true)
		if err != nil {
			return nil, err
		}
		if err := fStorage.readFileStorage(&ds); err != nil {
			return nil, err
		}
		if err := fStorage.close(); err != nil {
			return nil, err
		}
		logger.LogInfo("file storage data", zap.Int("count", len(ds.data)))
		return &ds, nil
	}

}

func (ds *DataStorage) Save(shortLink, originURL string) {
	if ds.fileStorage {
		fStorage, err := newFileStorage(false)
		if err != nil {
			logger.LogInfo("open file error", zap.Error(err))
		}
		fStorage.dataFill(shortLink, originURL)
		if err := fStorage.writeFileStorage(); err != nil {
			logger.LogInfo("write file error", zap.Error(err))
		}
		if err := fStorage.close(); err != nil {
			logger.LogInfo("close file error", zap.Error(err))
		}
	}
	ds.data[shortLink] = originURL
}

func (ds *DataStorage) Find(shortLink string) (string, error) {
	if originURL, ok := ds.data[shortLink]; ok {
		return originURL, nil
	}
	return "", errors.New("the object does not exist in storage")
}
