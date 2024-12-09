package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
)

type storager interface {
	InsertOneUrl(ctx context.Context, url *Url) error
	InsertMultipleUrl(ctx context.Context, urls *[]*Url) error
	SelectOneUrl(ctx context.Context, shortUrl string) (string, error)
	DeletePairUrl(ctx context.Context, shortUrl string) error
}

type inMemoryStorage struct {
	urlInfo map[string]string
	currentID uint
}

type inFileStorage struct {
	filename string
	currentID uint

}

type fileRD struct {
	file *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

type inDataBaseStorage struct {
	*sql.DB
}

type Url struct {
	ID uint
	CorrelationId string `json:"correlation_id,omitempty"`
	ShortUrl string `json:"short_url,omitempty"`
	OriginalUrl string `json:"original_url"`
}

type Storage struct {
	storager
}

type storageConfig struct {
	StorageType storageType
	Parameter string
}
