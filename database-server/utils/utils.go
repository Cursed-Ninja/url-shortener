package utils

import (
	"crypto/md5"
	b64 "encoding/base64"
	"time"
)

func KeyGenerationService(url string) string {
	hashedUrl := md5.New().Sum([]byte(url))
	encodedUrl := b64.StdEncoding.EncodeToString(hashedUrl)
	return encodedUrl
}

func GetExpirationTime() time.Time {
	currentTime := time.Now()
	expirationTime := currentTime.AddDate(0, 1, 0)
	return expirationTime
}
