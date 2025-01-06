package main

type ShortUrl struct {
	Id        int    `json:"id"`
	Url       string `json:"url"`
	ShortCode string `json:"short_code"`
}
