package middleware

import (
	"fmt"
	"net/http"
)

func CheckContentType(next http.Handler) http.Handler {
	const expectedHeader = "text/plain"

	return http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != expectedHeader {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)

		fmt.Println(r.Method)
	}))

}
