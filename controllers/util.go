package controllers

import (
	"math/rand"
	"time"
)

func randomInteger(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = NumberCharset[rand.Intn(len(NumberCharset))]
	}

	return string(b)
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = Charset[rand.Intn(len(Charset))]
	}

	return string(b)
}
