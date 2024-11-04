package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"

	config "github.com/PhenHF/url-shortener/internal/config"
	handler "github.com/PhenHF/url-shortener/internal/handler"
	middleware "github.com/PhenHF/url-shortener/internal/middleware"
	service "github.com/PhenHF/url-shortener/internal/service"
	storage "github.com/PhenHF/url-shortener/internal/storage"
)

func main() {
	var urlStorage = storage.UrlStorage{}

	rt := chi.NewRouter()

	rt.Use(middleware.CheckContentType)

	rt.Post(`/`, handler.ReturnShortUrl(service.GetShortUrl, &urlStorage))
	rt.Get(`/{id}`, handler.RedirectToOriginalUrl(&urlStorage))

	run(rt)
}

func run(rt *chi.Mux) {
	config.GetNetAddr()
	flag.Parse()

	err := http.ListenAndServe(config.NetAddress.StartServer, rt)
	if err != nil {
		panic(err)
	}
}
