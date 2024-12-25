package app

//go:generate mockgen -source=./shortener.go -destination=../mock/shortener.go -package=mock
type ShortenerHandler interface {
	RandomID() string
	CreateShortURL(body string) (string, error)
	GetShortURL(pathURL string) (string, error)
}
