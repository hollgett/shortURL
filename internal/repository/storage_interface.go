package repository

import "context"

type Storage interface {
	Save(shortLink, originURL string) error
	Find(shortLink string) (string, error)
	Close() error
	Ping(context.Context) error
}
