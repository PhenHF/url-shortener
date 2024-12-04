package storage

func idIncrement() uint {
	lenghtUrlStorage := len(UrlStorage)

	if lenghtUrlStorage < 1 {
		return 1
	}

	return UrlStorage[lenghtUrlStorage-1].ID + 1
}
