package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	middlewareLogger "github.com/PhenHF/url-shortener/internal/middleware/logger"
	storage "github.com/PhenHF/url-shortener/internal/storage"
	"go.uber.org/zap"
)

func RedirectToOriginalUrl() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("id") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		shortUrl := r.PathValue("id")
		res, err := storage.ReadyStorage.SelectOneUrl(r.Context(), shortUrl)
		if err != nil {
			if err != io.EOF {
				middlewareLogger.Log.Error("ERRROR", zap.String("msg", err.Error()))
			}
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Set("Location", res)
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}

func CreateShortUrl(generator func() string, resultAddr string) http.HandlerFunc {
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

		var myErr *storage.UniqueUrlError
		userID := r.Context().Value("user_id").(uint)
		if err := storage.ReadyStorage.InsertOneUrl(r.Context(), &url, userID); err != nil {
			if !errors.As(err, &myErr) {
				middlewareLogger.Log.Error("ERROR", zap.String("msg", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			shortUrl := strings.Split(err.Error(), ":")[1]
			url.ShortUrl = shortUrl
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

func CreateBatchShortUrl(generator func() string, resultAddr string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type resultUrl struct {
			CorrelationId string `json:"correlation_id"`
			ShortUrl      string `json:"short_url"`
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
		if err := storage.ReadyStorage.InsertMultipleUrl(r.Context(), &urls); err != nil {
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

func ReturnAllShortUrl(resultAddr string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id").(uint)

		urls, err := storage.ReadyStorage.SelectAllPairUrl(r.Context(), userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		type resultUrl struct {
			OriginalUrl string `json:"original_url"`
			ShortUrl    string `json:"short_url"`
		}

		resultUrls := make([]resultUrl, 0)

		for _, v := range *urls {
			resultUrls = append(resultUrls, resultUrl{OriginalUrl: v.OriginalUrl, ShortUrl: resultAddr + v.ShortUrl})
		}

		if len(resultUrls) < 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		response, err := json.Marshal(resultUrls)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write(response)

	})
}
