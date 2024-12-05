package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func readerConsole() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	url, err := reader.ReadString('\n')
	url = strings.TrimSpace(url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func clientShortener(url string) error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	fmt.Printf("original url: %#v\n", url)
	body := []byte(url)
	request, err := http.NewRequest(http.MethodPost, `http://localhost:8080/`, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "text/plain")

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	cnt, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Printf("status code: %v \t short url: %#v\n", response.StatusCode, string(cnt))

	request, err = http.NewRequest(http.MethodGet, string(cnt), nil)
	if err != nil {
		return err
	}
	response, err = client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	fmt.Printf("status code: %v\t\n", response.StatusCode)
	if _, err := io.Copy(io.Discard, response.Body); err != nil {
		return err
	}
	locationGet := response.Header.Get("Location")
	fmt.Printf("Location get: %s\r\n", locationGet)
	if url != locationGet {
		return errors.New("original url and short response not equal")
	}

	return nil
}

func main() {
	url, err := readerConsole()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := clientShortener(url); err != nil {
		fmt.Println(err)
		return
	}
}
