package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
)

func readerConsole() (string, error) {
	// reader := bufio.NewReader(os.Stdin)
	// url, err := reader.ReadString('\n')
	// url = strings.TrimSpace(url)
	// if err != nil {
	// 	return "", err
	// }
	url := "https://www.youtube.com/"
	return url, nil
}

type Client struct {
	httpClient *resty.Client
}

func NewClient() *Client {
	client := resty.NewWithClient(&http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}).
		SetBaseURL("http://localhost:8080").
		SetHeader("Content-Type", "text/plain; charset=utf-8")

	return &Client{
		httpClient: client,
	}
}

func clientShortener(originalURL string) error {
	fmt.Printf("original url: %#v\n", originalURL)
	client := NewClient()
	resp, err := client.httpClient.R().SetBody(originalURL).Post("/")
	if err != nil {
		return err
	}
	resShort := resp.String()
	fmt.Printf("status code post: %v \t short url: %#v\n", resp.StatusCode(), resShort)
	shortURL, err := url.Parse(resShort)
	if err != nil {
		return err
	}
	fmt.Println(shortURL.Path)
	resp, err = client.httpClient.R().Get(shortURL.Path)
	if err != nil {
		return err
	}
	fmt.Printf("status code get: %v\t\n", resp.StatusCode())
	locationGet := resp.Header().Get("Location")
	fmt.Printf("Location get: %s\r\n", locationGet)
	if originalURL != locationGet {
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
