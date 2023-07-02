package random

import (
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func NewRandomString(size int) string {
	random := rand.NewSource(time.Now().UnixNano())
	b := make([]rune, size)
	for i := range b {
		b[i] = letterRunes[random.Int63()%int64(len(letterRunes))]
	}
	return string(b)
}
