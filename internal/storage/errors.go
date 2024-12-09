package storage

import "fmt"

type UniqueUrlError struct {
	shortUrl string
}

func (uue *UniqueUrlError) Error() string {
	return fmt.Sprintf("url already exists:%s", uue.shortUrl)
}

func newUrlError(shortUrl string) error {
	return &UniqueUrlError{
		shortUrl: shortUrl,
	}
}
