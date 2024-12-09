package config

import (
	"flag"

	"github.com/PhenHF/url-shortener/internal/storage"
)

func init() {
	loadFlags()
	loadEnv()
}

func loadFlags() {
	var dbParam string
	var fileParam string

	flag.StringVar(&NetAddress.StartServer, "a", ":8080", "addr for start a server")
	flag.StringVar(&NetAddress.ResultAddr, "b", "http://localhost:8080/", "addr for base result URL")

	flag.StringVar(&dbParam, "d", "", "username:password for connection to DB")
	flag.StringVar(&fileParam, "f", "", "filepath for save url in json file")
	flag.Parse()
	
	if dbParam != "" {
		StorageConfig.StorageType = storage.InDataBase
		StorageConfig.Parameter = dbParam
		return
	}

	if fileParam != "" {
		StorageConfig.StorageType = storage.InFile
		StorageConfig.Parameter = fileParam
	}
}
