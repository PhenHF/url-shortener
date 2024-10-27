package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/PhenHF/url-shortener/internal/app"
)


func main() {
	http.HandleFunc(`/`, Handler)
	http.HandleFunc(`/{id}`, Handler)
	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}


func Handler(w http.ResponseWriter, r *http.Request) {
	const expectedHeader = "text/plain"
	
	// fmt.Println(app.OriginalShortURL)
	// fmt.Println(app.ShortOriginalURL)

	
	if ct := r.Header.Get("Content-Type"); ct != expectedHeader {
		fmt.Println("1")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodGet {
		fmt.Println(app.ShortOriginalURL)
		if r.PathValue("id") == "" {
			fmt.Println("2")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		originUrl, ok := app.ShortOriginalURL[r.PathValue("id")]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", originUrl)
		// w.WriteHeader(http.StatusTemporaryRedirect)		
		// w.Write([]byte(originUrl))
	
	} else if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
				
		w.Header().Set("Content-Type", "text/plain")
		
		shortUrl := app.GetShortUrl()
		
		app.ShortOriginalURL[shortUrl] = string(body)
		w.WriteHeader(http.StatusCreated)		
		w.Write([]byte("http://localhost:8080/" + shortUrl))
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}