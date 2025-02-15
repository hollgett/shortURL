package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func clientShortenerJSON(originalURL string) error {

	client := NewClient()
	req := []RequestBatch{
		RequestBatch{
			Id:       "1",
			Original: originalURL,
		},
		RequestBatch{
			Id:       "2",
			Original: originalURL,
		},
		RequestBatch{
			Id:       "3",
			Original: originalURL,
		},
		RequestBatch{
			Id:       "4",
			Original: originalURL,
		},
	}
	fmt.Printf("original url: %#v\n", req)
	request, err := json.Marshal(req)
	if err != nil {
		return err
	}
	var bw bytes.Buffer
	gw := gzip.NewWriter(&bw)
	_, err = gw.Write(request)
	if err != nil {
		return fmt.Errorf("gzip writer: %w", err)
	}
	gw.Close()
	resp, err := client.httpClient.R().SetBody(bw.Bytes()).
		SetHeader("Content-Type", "application/x-gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Encoding", "gzip").
		Post("/api/shorten/batch")
	if err != nil {
		return fmt.Errorf("request post: %w", err)
	}
	fmt.Println(resp.Header().Get("Content-Encoding"))
	fmt.Println(resp.String())
	gr, err := gzip.NewReader(resp.RawBody())
	if err != nil {
		return fmt.Errorf("gzip reader: %w", err)
	}
	defer gr.Close()
	var br bytes.Buffer
	_, err = io.Copy(&br, gr)
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	response := []ResponseBatch{}
	json.Unmarshal(br.Bytes(), &response)
	fmt.Printf("status code post: %v \t short url: %#v\n", resp.StatusCode(), response)
	fmt.Printf("%+v\n", response)
	// shortURL, err := url.Parse(response.ResponseURL)
	// if err != nil {
	// 	return fmt.Errorf("parse url: %w", err)
	// }
	// fmt.Println(shortURL.Path)
	// resp, err = client.httpClient.R().Get(shortURL.Path)
	// if err != nil {
	// 	return fmt.Errorf("request get: %w", err)
	// }
	// fmt.Printf("status code get: %v\t\n", resp.StatusCode())
	// locationGet := resp.Header().Get("Location")
	// fmt.Printf("Location get: %s\r\n", locationGet)
	// if originalURL != locationGet {
	// 	return errors.New("original url and short response not equal")
	// }
	return nil
}

func clientShortenerGzip(originalURL string) error {
	fmt.Printf("original url: %#v\n", originalURL)
	client := NewClient()
	var bw bytes.Buffer
	gw := gzip.NewWriter(&bw)
	_, err := gw.Write([]byte(originalURL))
	if err != nil {
		return fmt.Errorf("gzip writer: %w", err)
	}
	fmt.Println("bytes", bw.String())
	gw.Close()
	resp, err := client.httpClient.R().SetBody(bw.Bytes()).
		SetHeader("Content-Type", "application/x-gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Encoding", "gzip").
		Post("/")
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

	if err := clientShortenerJSON(url); err != nil {
		fmt.Println(err)
		return
	}
}
