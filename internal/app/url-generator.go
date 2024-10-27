package app

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var lenght = 8

func GetShortUrl() string {
	var su []rune
	for i := 0; i < lenght; i++ {
		su = append(su, letterRunes[rand.Intn(len(letterRunes))])

	}
	return string(su)
}
