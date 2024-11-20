package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	config "github.com/PhenHF/url-shortener/internal/config"
	handler "github.com/PhenHF/url-shortener/internal/handler"
	middlewareCopmpress "github.com/PhenHF/url-shortener/internal/middleware/httpcompress"
	middlewareLogger "github.com/PhenHF/url-shortener/internal/middleware/logger"
	service "github.com/PhenHF/url-shortener/internal/service"
	storage "github.com/PhenHF/url-shortener/internal/storage"
)

func main() {
	var urlStorage = storage.UrlStorage{}
	rt := chi.NewRouter()
	rt.Use(middlewareLogger.RequestLogger)
	rt.Use(middlewareCopmpress.GzipMiddleware)
	rt.Post(`/api/shorten`, handler.ReturnShortUrl(service.GetShortUrl, &urlStorage, config.NetAddress.ResultAddr))
	rt.Get(`/{id}`, handler.RedirectToOriginalUrl(&urlStorage))
	run(rt)
}

func run(rt *chi.Mux) error {
	if err := middlewareLogger.Initialize("INFO"); err != nil {
		return err
	}

	err := http.ListenAndServe(config.NetAddress.StartServer, rt)
	if err != nil {
		panic(err)
	}
	return nil
}
