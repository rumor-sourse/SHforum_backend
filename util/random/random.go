package random

import (
	"math/rand"
	"time"
)

const (
	Charset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	NumberCharset = "0123456789"
)

func IntegerN(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = NumberCharset[rand.Intn(len(NumberCharset))]
	}

	return string(b)
}

func StringN(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = Charset[rand.Intn(len(Charset))]
	}

	return string(b)
}
