package model

type ShortURL struct {
	ID          string `json:"id"`
	FullURL     string `json:"fullURL"`
	ShortCode   string `json:"shortCode"`
	AccessCount int    `json:"accessCount"`
}
