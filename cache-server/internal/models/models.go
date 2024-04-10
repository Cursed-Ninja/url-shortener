package models

import "time"

type RequestModel struct {
	Url       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ShortenRequestModel struct {
	Url       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ShortenResponseModel struct {
	ShortUrlPath string `json:"shorturlpath"`
}

type RedirectRequestModel struct {
	ShortUrlPath string `json:"shorturlpath"`
}

type RedirectResponseModel struct {
	Url string `json:"redirecturl"`
}

type ResponseModel struct {
	Url string `json:"url"`
}
