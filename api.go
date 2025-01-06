package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type api struct {
	addr string
}

var idCounter int
var shortendedUrls = map[int]ShortUrl{}
var shortCodeIndex = map[string]int{}

func redirectShortUrlToDestinationHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")
	shortUrlId, ok := shortCodeIndex[shortCode]
	if !ok {
		http.Error(w, "Short url not found", http.StatusNotFound)
		return
	}
	shortUrl, ok := shortendedUrls[shortUrlId]
	if !ok {
		http.Error(w, "Failed to retrieve ShortUrl", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, shortUrl.Url, http.StatusSeeOther)
}

func (a *api) createShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var url Url
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortUrl, err := createAndInsertShortUrl(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(shortUrl)
	w.WriteHeader(http.StatusCreated)
}

func createAndInsertShortUrl(u Url) (*ShortUrl, error) {
	if u.Url == "" {
		return nil, errors.New("url is missing")
	}
	if u.ShortCode == "" {
		return nil, errors.New("short_code is missing")
	}
	if _, ok := shortCodeIndex[u.ShortCode]; ok {
		return nil, fmt.Errorf("short code '%s' already exists", u.ShortCode)
	}

	id := idCounter
	idCounter += 1
	shortUrl := ShortUrl{
		Id:        id,
		ShortCode: u.ShortCode,
		Url:       addHTTP(u.Url),
	}
	shortendedUrls[shortUrl.Id] = shortUrl
	shortCodeIndex[shortUrl.ShortCode] = shortUrl.Id
	return &shortUrl, nil
}

func addHTTP(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}
	return url
}
