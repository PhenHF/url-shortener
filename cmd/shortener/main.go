package main

import (
	"net/http"

	"github.com/PhenHF/url-shortener/internal/app"
)

func init() {
	http.HandleFunc(`/`, app.ReturnShortUrl)
	http.HandleFunc(`/{id}`, app.RedirectToOriginalUrl)
}

func main() {
	run()
}

func run() {
	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}
