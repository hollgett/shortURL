package app

//go:generate mockgen -source=./shortener_interface.go -destination=../mocks/shortener.go -package=mocks
type ShortenerHandler interface {
	RandomID() string
	CreateShortURL(requestData string) (string, error)
	GetShortURL(pathURL string) (string, error)
}
