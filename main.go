package main

import (
	"log"
	"net/http"
)

// ANSI escape codes for colors
const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

func main() {
	log.Printf(Green + "starting short url service..." + Reset)
	api := &api{addr: ":8080"}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{shortCode}", redirectShortUrlToDestinationHandler)
	mux.HandleFunc("POST /shorten", api.createShortUrlHandler)
	srv := &http.Server{
		Addr:    api.addr,
		Handler: mux,
	}
	log.Fatal(srv.ListenAndServe())
}
