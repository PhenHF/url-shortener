package storage

type storageType int

const (
	InMemory storageType = iota
	InDataBase
	InFile
)

func BuildDB(storageConfig storageConfig) Storage {
	switch storageConfig.StorageType {
	case InDataBase:
		db := newInDataBaseStorage(storageConfig)
		db.initDB()
		return Storage{db}
	case InFile:
		fs := newFileStorage(storageConfig.Parameter)
		return Storage{fs}
	default:
		ms, err := newInMemoryStorage()
		if err != nil {
			panic(err)
		}
		return Storage{ms}
	}
}

func NewStorageConfig() *storageConfig {
	return &storageConfig{
		StorageType: 0,
		Parameter: "",
	}
}