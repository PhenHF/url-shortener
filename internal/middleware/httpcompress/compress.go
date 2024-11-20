package httpcompress

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

// Middleware for commpress and decompress request and response
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		if !slices.Contains(contentType, r.Header.Get("Content-Type")) {
			fmt.Println("blla")
			next.ServeHTTP(w, r)
			return
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fmt.Println("blla")
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}
		
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}
		next.ServeHTTP(ow, r)
	})
}