package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	handler "github.com/PhenHF/url-shortener/internal/handler"
	middleware "github.com/PhenHF/url-shortener/internal/middleware"
)


func main() {
	rt := chi.NewRouter()

	rt.Use(middleware.CheckContentType)
	
	rt.Post(`/`, handler.ReturnShortUrl)
	rt.Get(`/{id}`, handler.RedirectToOriginalUrl)

	run(rt)
}

func run(rt *chi.Mux) {
	err := http.ListenAndServe(`:8080`, rt)
	if err != nil {
		panic(err)
	}

	
}
