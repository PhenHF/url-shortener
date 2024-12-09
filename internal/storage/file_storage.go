package storage

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
)

func newFileStorage(filename string) *inFileStorage {
	return &inFileStorage{
		filename:  filename,
		currentID: 0,
	}
}

func newFileRD(filename string) (*fileRD, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &fileRD{
		file:    file,
		decoder: json.NewDecoder(file),
		encoder: json.NewEncoder(file),
	}, nil
}

func (fs *inFileStorage) InsertOneUrl(ctx context.Context, url *Url) error {
	frd, err := newFileRD(fs.filename)
	if err != nil {
		panic(err)
	}
	defer frd.file.Close()

	fs.inMemoryStorageIdIncrement()
	url.ID = fs.currentID
	return frd.encoder.Encode(url)
}

func (fs *inFileStorage) InsertMultipleUrl(ctx context.Context, urls *[]*Url) error {
	frd, err := newFileRD(fs.filename)
	if err != nil {
		panic(err)
	}
	defer frd.file.Close()

	for _, v := range *urls {
		fs.inMemoryStorageIdIncrement()
		v.ID = fs.currentID
		err := frd.encoder.Encode(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fs *inFileStorage) SelectOneUrl(ctx context.Context, shortUrl string) (string, error) {
	frd, err := newFileRD(fs.filename)
	if err != nil {
		panic(err)
	}
	defer frd.file.Close()

	for {
		readUrl, err := fs.Read(frd)

		if err != nil {
			if err == io.EOF {
				return "", err
			}

			return "", err
		}
		if readUrl.CorrelationId == shortUrl {
			return readUrl.OriginalUrl, nil
		}
	}
}

func (fs *inFileStorage) Read(frd *fileRD) (*Url, error) {
	url := &Url{}
	if err := frd.decoder.Decode(&url); err != nil {
		return nil, err
	}
	return url, nil
}

func (fs *inFileStorage) DeletePairUrl(ctx context.Context, shortUrl string) error {
	return errors.New("method is not allowed")
}

func (fs *inFileStorage) inMemoryStorageIdIncrement() {
	fs.currentID++
}
