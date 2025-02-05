package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/hollgett/shortURL.git/internal/logger"
	"github.com/hollgett/shortURL.git/internal/models"
	"go.uber.org/zap"
)

type fileStorage struct {
	file *os.File
	data map[string]string
}

func NewFileStorage(src string) (Storage, error) {
	file, err := os.OpenFile(src, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	fs := &fileStorage{
		file: file,
		data: make(map[string]string),
	}
	fileData, err := fs.load()
	if err != nil {
		return nil, err
	}
	for _, v := range fileData {
		fs.data[v.Short] = v.Original
	}
	logger.LogInfo("file restore", zap.Int("len", len(fileData)))
	return fs, nil
}

func (fs *fileStorage) Close() error {
	return fs.file.Close()
}

func (fs *fileStorage) load() ([]models.FileStorageData, error) {
	scan := bufio.NewScanner(fs.file)
	fileData := []models.FileStorageData{}
	data := models.FileStorageData{}

	for scan.Scan() {
		if err := json.Unmarshal(scan.Bytes(), &data); err != nil {
			return nil, fmt.Errorf("decode data: %w", err)
		}
		fileData = append(fileData, data)
	}
	if err := scan.Err(); err != nil {
		return nil, fmt.Errorf("scanner data: %w", err)
	}
	return fileData, nil
}

func (fs *fileStorage) write(shortLink, originURL string) error {
	dataSave := models.FileStorageData{
		Short:    shortLink,
		Original: originURL,
	}
	data, err := json.Marshal(dataSave)
	if err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	data = append(data, '\n')
	if _, err := fs.file.Write(data); err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	if err = fs.file.Sync(); err != nil {
		return fmt.Errorf("sync error: %w", err)
	}
	return nil
}

func (fs *fileStorage) Save(shortLink, originURL string) error {
	fs.data[shortLink] = originURL
	if err := fs.write(shortLink, originURL); err != nil {
		return err
	}
	return nil
}
func (fs *fileStorage) Find(shortLink string) (string, error) {
	if originURL, ok := fs.data[shortLink]; ok {
		return originURL, nil
	}
	return "", errors.New("the object does not exist in storage")
}

func (fs *fileStorage) Ping(context.Context) error {
	if _, err := fs.file.Stat(); err != nil {
		return err
	}
	return nil
}
