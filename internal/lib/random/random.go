package random

import (
	"math/rand"
	"time"
)

func NewRandomString(length int) string {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune(
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
			"abcdefghijklmnopqrstuvwxyz" +
			"0123456789")

	symbols := make([]rune, length)
	for i := range symbols {
		symbols[i] = chars[rand.Intn(len(chars))]
	}

	return string(symbols)
}
