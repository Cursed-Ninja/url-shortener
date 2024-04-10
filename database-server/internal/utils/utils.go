package utils

import (
	"crypto/md5"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const (
	base         int    = 62
	characterSet string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func toBase62(b int) string {
	encoded := ""
	for b > 0 {
		r := b % base
		b /= base
		encoded = string(characterSet[r]) + encoded
	}
	return encoded
}

func KeyGenerationService(url string) string {
	hashedUrl := md5.New().Sum([]byte(url))
	encodedUrl := ""
	for _, b := range hashedUrl {
		encodedUrl += toBase62(int(b))
	}

	shortUrl := ""

	n := len(encodedUrl)

	for len(shortUrl) < 7 {
		shortUrl += string(encodedUrl[rand.Intn(n)])
	}

	return shortUrl
}

func GetExpirationTime(expiryTime time.Time) time.Time {
	currentTime := time.Now()

	if currentTime.After(expiryTime) {
		return currentTime.AddDate(0, 1, 0)
	}

	return expiryTime
}

func GenerateRequestId() string {
	return uuid.New().String()
}
