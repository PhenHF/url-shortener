package storage

import (
	"errors"
)

type UrlStorage struct {
	shortOriginalURL map[string]string
}

func (us *UrlStorage) Add(baseUrl, shortUrl string) {
	if us.shortOriginalURL == nil {
		us.shortOriginalURL = make(map[string]string)
	}

	us.shortOriginalURL[shortUrl] = baseUrl
}

func (us UrlStorage) Get(shortUrl string) (string, error) {
	baseUrl, ok := us.shortOriginalURL[shortUrl]
	if !ok {
		return "", errors.New("no url for this ID")
	}

	return baseUrl, nil
}
