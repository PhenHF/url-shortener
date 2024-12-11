package storage

import (
	"math/rand"
	"time"
)

var ReadyStorage Storage

type storageType int

const (
	InMemory storageType = iota
	InDataBase
	InFile
)

func BuildDB(storageConfig storageConfig) {
	switch storageConfig.StorageType {
	case InDataBase:
		db := newInDataBaseStorage(storageConfig)
		db.initDB()
		ReadyStorage = Storage{db}
		// case InFile:
		// 	fs := newFileStorage(storageConfig.Parameter)
		// 	return Storage{fs}
		// default:
		// 	ms, err := newInMemoryStorage()
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	return Storage{ms}
		// }
	}
}

func NewStorageConfig() *storageConfig {
	return &storageConfig{
		StorageType: 0,
		Parameter:   "",
	}
}

func generateShortUrl() string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var lenght = 8

	var su []rune
	for i := 0; i < lenght; i++ {
		su = append(su, letterRunes[rand.Intn(len(letterRunes))])

	}
	return string(su)
}
