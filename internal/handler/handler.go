package handler

import (
	"fmt"
	"io"
	"net/http"

	storage "github.com/PhenHF/url-shortener/internal/storage"
)

func RedirectToOriginalUrl(urlStorage *storage.UrlStorage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("id") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		originUrl, err := urlStorage.Get(r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", originUrl)

		w.WriteHeader(http.StatusTemporaryRedirect)

	})
}

func ReturnShortUrl(generator func() string, urlStorage *storage.UrlStorage, resultAddr string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}

		if len(body) == 0 {
			w.WriteHeader(http.StatusBadRequest)
		}

		short := generator()
		urlStorage.Add(string(body), short)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(resultAddr + short))
	})
}
