package storage

import (
	"context"
	"fmt"
)

func newInMemoryStorage() (*inMemoryStorage, error) {
	return &inMemoryStorage{
		urlInfo:   make(map[string]string),
		currentID: 0,
	}, nil
}

func (ms *inMemoryStorage) InsertOneUrl(ctx context.Context, url *Url) error {
	ms.inMemoryStorageIdIncrement()
	url.ID = ms.currentID
	ms.urlInfo[url.CorrelationId] = url.OriginalUrl
	return nil
}

func (ms *inMemoryStorage) InsertMultipleUrl(ctx context.Context, urls *[]*Url) error {
	for _, v := range *urls {
		fmt.Println(v.CorrelationId)
		ms.inMemoryStorageIdIncrement()
		v.ID = ms.currentID
		ms.urlInfo[v.CorrelationId] = v.OriginalUrl
	}
	return nil
}

func (ms *inMemoryStorage) DeletePairUrl(ctx context.Context, shortUrl string) error {
	delete(ms.urlInfo, shortUrl)
	return nil
}

func (ms inMemoryStorage) SelectOneUrl(ctx context.Context, shortUrl string) (string, error) {
	originalUrl, ok := ms.urlInfo[shortUrl]
	if !ok {
		return "", fmt.Errorf("no url for shortUrl equal %s", shortUrl)
	}
	return originalUrl, nil
}

func (ms *inMemoryStorage) inMemoryStorageIdIncrement() {
	ms.currentID++
}
