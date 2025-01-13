package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hollgett/shortURL.git/internal/config"
)

type fileStorage struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

func readFileStorage(storage *map[string]string) error {
	file, err := os.OpenFile(config.Cfg.FileStorage, os.O_RDONLY|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("open file read: %w", err)
	}
	reader := bufio.NewScanner(file)

	for reader.Scan() {
		var fStorage fileStorage
		if err := json.Unmarshal(reader.Bytes(), &fStorage); err != nil {
			return fmt.Errorf("decode data: %w", err)
		}
		(*storage)[fStorage.Short] = fStorage.Original
	}
	if err := reader.Err(); err != nil {
		return fmt.Errorf("scanner data: %w", err)
	}
	return nil
}

func writeFileStorage(shortLink, originURL string) error {
	file, err := os.OpenFile(config.Cfg.FileStorage, os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("open file write: %w", err)
	}
	data, err := json.Marshal(fileStorage{
		Short:    shortLink,
		Original: originURL,
	})
	if err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	data = append(data, '\n')
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	return nil
}
