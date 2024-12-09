package config

import (
	"os"

	"github.com/PhenHF/url-shortener/internal/storage"
)

func loadEnv() {
	if serverAddr := os.Getenv("SERVER_ADDRESS"); serverAddr != "" {
		NetAddress.StartServer = serverAddr
	}
	
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		NetAddress.ResultAddr = baseURL
	}
	
	if postgresCred := os.Getenv("DATABASE_DSN"); postgresCred != "" {
		StorageConfig.Parameter = postgresCred
		StorageConfig.StorageType = storage.InDataBase
		return
	}

	if jsonFilename := os.Getenv("FILE_STORAGE_PATH"); jsonFilename != "" && StorageConfig.StorageType != storage.InDataBase{
		StorageConfig.Parameter = jsonFilename
		StorageConfig.StorageType = storage.InFile
		return
	}

	if StorageConfig.Parameter == "" && StorageConfig.StorageType != storage.InDataBase && StorageConfig.StorageType != storage.InFile {
		StorageConfig.StorageType = storage.InMemory
	}
}