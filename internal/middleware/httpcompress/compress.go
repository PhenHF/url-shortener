package httpcompress

import (
	"net/http"
	"slices"
	"strings"
)

// Middleware for write and read commpress request and response
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		if !slices.Contains(contentType, r.Header.Get("Content-Type")) {
			next.ServeHTTP(w, r)
			return
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
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
