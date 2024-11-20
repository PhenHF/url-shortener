package storage

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func init() {
	us, err := NewUrlConsumer("url.json")
	if err != nil {
		log.Printf("can't read url.json with err:%v", err)
	}
	defer us.Close()
	us.ReadAll()
}

type Url struct {
	ID uint `json:"uuid,omitempty"`
	ShortUrl string `json:"short_url,omitempty"`
	OriginalUrl string `json:"original_url"`
}

var UrlStorage []Url

type UrlProducer struct {
	file *os.File
	encoder *json.Encoder
}

func NewUrlProducer(filename string) (*UrlProducer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &UrlProducer{
		file: file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (up *UrlProducer) WriteUrl(url *Url) error {
	url.ID = idIncrement()
	UrlStorage = append(UrlStorage, *url)
	return up.encoder.Encode(url)
}

func (up *UrlProducer) Close() {
	up.file.Close()
}

type UrlConsumer struct {
	file *os.File
	decoder *json.Decoder
}

func NewUrlConsumer(filename string) (*UrlConsumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &UrlConsumer{
		file: file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (uc *UrlConsumer) Read() (*Url, error) {
	url := &Url{}
	if err := uc.decoder.Decode(&url); err != nil {
		return nil, err
	}
	return url, nil
}
 
func (uc *UrlConsumer) ReadAll() (error) {
	for {
		readUrl, err := uc.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			
			return err
		}
		UrlStorage = append(UrlStorage, *readUrl)
	}
}

func (uc *UrlConsumer) Close() {
	uc.file.Close()
}