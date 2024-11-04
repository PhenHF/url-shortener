package handler

import (
	"fmt"
	"io"
	"net/http"

	servise "github.com/PhenHF/url-shortener/internal/service"
	storage "github.com/PhenHF/url-shortener/internal/storage"
)

// func checkContentType(w http.ResponseWriter, r *http.Request) {
// 	const expectedHeader = "text/plain"
// 	if ct := r.Header.Get("Content-Type"); ct != expectedHeader {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// }

func RedirectToOriginalUrl(w http.ResponseWriter, r *http.Request) {
	if r.PathValue("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originUrl, ok := storage.ShortOriginalURL[r.PathValue("id")]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originUrl)

	w.WriteHeader(http.StatusTemporaryRedirect)
}

func ReturnShortUrl(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
	}

	shortUrl := servise.GetShortUrl()

	storage.ShortOriginalURL[shortUrl] = string(body)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + shortUrl))
}
