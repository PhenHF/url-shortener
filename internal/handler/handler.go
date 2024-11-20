package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	storage "github.com/PhenHF/url-shortener/internal/storage"
)

func RedirectToOriginalUrl(urlStorage *[]storage.Url) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("id") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		shortUrl := r.PathValue("id")
		fmt.Println(shortUrl)
		for _, url := range *urlStorage {
			if url.ShortUrl == shortUrl {
				w.Header().Set("Location", url.OriginalUrl)
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
		}
		w.WriteHeader(http.StatusBadRequest)
	})
}

func ReturnShortUrl(generator func() string, resultAddr string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resultUrl struct {
			Url string `json:"response"`
		}
		
		url := storage.Url{}

		urlProducer, err := storage.NewUrlProducer("url.json")
		if err != nil {
			return
		}
		defer urlProducer.Close()
		
		
		if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
			return 
		}
		
		if len(url.OriginalUrl) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		url.ShortUrl = generator()
		
		urlProducer.WriteUrl(&url)
		resultUrl.Url = resultAddr + url.ShortUrl

		response, err := json.Marshal(resultUrl)
		if err != nil {
			return
		}


		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	})
}
