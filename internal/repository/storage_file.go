package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hollgett/shortURL.git/internal/config"
)

type fileStorageData struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

type fileStorage struct {
	file *os.File
	data fileStorageData
}

func newFileStorage(readF bool) (*fileStorage, error) {
	var flag int
	switch readF {
	case true:
		flag = os.O_RDONLY | os.O_CREATE
	case false:
		flag = os.O_WRONLY | os.O_APPEND
	}
	file, err := os.OpenFile(config.Cfg.FileStorage, flag, 0666)
	if err != nil {
		return nil, err
	}
	return &fileStorage{
		file: file,
	}, nil
}

func (fs *fileStorage) close() error {
	return fs.file.Close()
}

func (fs *fileStorage) readFileStorage(storage *DataStorage) error {
	storage.fileStorage = true
	scan := bufio.NewScanner(fs.file)
	for scan.Scan() {
		if err := json.Unmarshal(scan.Bytes(), &fs.data); err != nil {
			return fmt.Errorf("decode data: %w", err)
		}
		storage.data[fs.data.Short] = fs.data.Original
	}
	if err := scan.Err(); err != nil {
		return fmt.Errorf("scanner data: %w", err)
	}
	return nil
}

func (fs *fileStorage) writeFileStorage() error {
	data, err := json.Marshal(fs.data)
	if err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	data = append(data, '\n')
	if _, err := fs.file.Write(data); err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	return nil
}

func (fs *fileStorage) dataFill(shLink, origURL string) {
	fs.data.Short = shLink
	fs.data.Original = origURL
}
