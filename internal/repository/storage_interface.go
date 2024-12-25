package repository

type Storage interface {
	Save(shortLink, originURL string)
	Find(shortLin string) (string, error)
}
