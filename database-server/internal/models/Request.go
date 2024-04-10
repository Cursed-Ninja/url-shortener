package models

import "time"

type RequestModel struct {
	Url string `json:"redirectUrl"`
}

type URL struct {
	Url          string
	ShortenedUrl string
	ExpiresAt    time.Time
}
