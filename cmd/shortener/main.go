package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	config "github.com/PhenHF/url-shortener/internal/config"
	handler "github.com/PhenHF/url-shortener/internal/handler"
	middlewareAuth "github.com/PhenHF/url-shortener/internal/middleware/auth"
	middlewareCopmpress "github.com/PhenHF/url-shortener/internal/middleware/httpcompress"
	middlewareLogger "github.com/PhenHF/url-shortener/internal/middleware/logger"
	service "github.com/PhenHF/url-shortener/internal/service"
	storage "github.com/PhenHF/url-shortener/internal/storage"
)

func main() {
	storage.BuildDB(*config.StorageConfig)
	rt := chi.NewRouter()
	rt.Use(middlewareLogger.RequestLogger)
	rt.Use(middlewareCopmpress.GzipMiddleware)
	rt.Use(middlewareAuth.AuthMiddleware)
	rt.Post(`/api/shorten`, handler.CreateShortUrl(service.GetShortUrl, config.NetAddress.ResultAddr))
	rt.Post(`/api/shorten/batch`, handler.CreateBatchShortUrl(service.GetShortUrl, config.NetAddress.ResultAddr))
	rt.Get(`/{id}`, handler.RedirectToOriginalUrl())
	rt.Get(`/api/user/urls`, handler.ReturnAllShortUrl(config.NetAddress.ResultAddr))
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
