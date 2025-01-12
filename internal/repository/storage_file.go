package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/hollgett/shortURL.git/internal/config"
	"github.com/hollgett/shortURL.git/internal/logger"
	"go.uber.org/zap"
)

type fileStorage struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

func readFileStorage(storage *map[string]string) error {
	path, err := pathToTemp()
	if err != nil {
		return fmt.Errorf("path file: %w", err)
	}
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0555)
	defer func() {
		if err := file.Close(); err != nil {
			logger.LogInfo("close file read", zap.Error(err))
		}
	}()
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	reader := bufio.NewReader(file)
	var fStorage fileStorage
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("read file: %w", err)
		}
		if len(data) == 0 && err == io.EOF {
			return nil
		}
		if err := json.Unmarshal(data, &fStorage); err != nil {
			return fmt.Errorf("decode data: %w", err)
		}
		(*storage)[fStorage.Short] = fStorage.Original
	}
}

func writeFileStorage(shortLink, originURL string) error {
	path, err := pathToTemp()
	if err != nil {
		return fmt.Errorf("path file: %w", err)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0333)
	defer func() {
		if err := file.Close(); err != nil {
			logger.LogInfo("close file write", zap.Error(err))
		}
	}()
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	data, err := json.Marshal(fileStorage{
		Short:    shortLink,
		Original: originURL,
	})
	if err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	writer := bufio.NewWriter(file)
	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("write data: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("write byte: %w", err)
	}

	return writer.Flush()
}

func pathToTemp() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get path: %w", err)
	}
	path := path.Join(dir, config.Cfg.FileStorage)
	logger.LogInfo("temp", zap.String("path", path))
	return path, nil
}
