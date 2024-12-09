package handler

import (
	"encoding/json"
	"io"
	"net/http"

	middlewareLogger "github.com/PhenHF/url-shortener/internal/middleware/logger"
	storage "github.com/PhenHF/url-shortener/internal/storage"
	"go.uber.org/zap"
)

func RedirectToOriginalUrl(urlStorage storage.Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("id") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		shortUrl := r.PathValue("id")
		res, err := urlStorage.SelectOneUrl(r.Context(), shortUrl)
		if err != nil {
			if err != io.EOF{
				middlewareLogger.Log.Error("ERRROR", zap.String("msg", err.Error()))
			}
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Set("Location", res)
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}

func ReturnShortUrl(generator func() string, resultAddr string, urlStorage storage.Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resultUrl struct {
			Url string `json:"response"`
		}

		url := storage.Url{}

		if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
			return
		}

		if len(url.OriginalUrl) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		url.ShortUrl = generator()
		
		err := urlStorage.InsertOneUrl(r.Context(), &url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resultUrl.Url = resultAddr + url.ShortUrl

		response, err := json.Marshal(resultUrl)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	})
}

func ReturnBatchShortUrl(generator func() string, resultAddr string, urlStorage storage.Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type resultUrl struct{
			CorrelationId string`json:"correlation_id"`
			ShortUrl string `json:"short_url"`
		}
		
		urls := []*storage.Url{}
		if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		resultUrls := make([]resultUrl, 0)

		for _, v := range urls {
			v.ShortUrl = generator()
			resultUrls = append(resultUrls, resultUrl{CorrelationId: v.CorrelationId, ShortUrl: resultAddr + v.ShortUrl})
		}
		
		if err := urlStorage.InsertMultipleUrl(r.Context(), &urls); err != nil {
			middlewareLogger.Log.Error("ERROR", zap.String("msg:", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(resultUrls)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)

	})
}