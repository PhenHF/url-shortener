package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/PhenHF/url-shortener/internal/app"
)


func main() {
	rt := chi.NewRouter()
	
	rt.Use(app.CheckContentType)
	
	rt.Post(`/`, app.ReturnShortUrl)
	rt.Get(`/{id}`, app.RedirectToOriginalUrl)

	run(rt)
}

func run(rt *chi.Mux) {
	err := http.ListenAndServe(`:8080`, rt)
	if err != nil {
		panic(err)
	}

	
}
