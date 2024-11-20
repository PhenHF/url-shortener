package handler

import (
	"encoding/json"
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
		var longUrl struct {
			Url string `json:"url"`
		}

		var resultUrl struct {
			Url string `json:"result"`
		}

		if err := json.NewDecoder(r.Body).Decode(&longUrl); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(longUrl.Url) == 0 {
			w.WriteHeader(http.StatusBadRequest)
		}

		short := generator()
		urlStorage.Add(longUrl.Url, short)

		resultUrl.Url = resultAddr + short
		
		response, err := json.Marshal(resultUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	})
}
