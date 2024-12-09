package config

import "github.com/PhenHF/url-shortener/internal/storage"

var NetAddress struct {
	StartServer string
	ResultAddr  string
}

var StorageConfig = storage.NewStorageConfig()
