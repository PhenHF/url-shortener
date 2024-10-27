package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhenHF/url-shortener/internal/app"
	"github.com/stretchr/testify/assert"
)


func TestRedirectToOriginalUrl(t *testing.T) {
	shortUrl := "stbfg"
	
	app.ShortOriginalURL[shortUrl] = "https://www.google.com/" 
	
	type want struct {
		code int
		location string		
	}
	
	tests := []struct {
		name string
		shorturl string
		contentType string
		want want
	}{
		{
			name: "test with correct shortUrl and content-type: text/plain",
			shorturl: shortUrl,
			contentType: "text/plain",
			want: want{
				code: 307,
				location: app.ShortOriginalURL[shortUrl], 
			},
		},

		{
			name: "test with NOT correct shortUrl and content-type: text/plain",
			shorturl: "asdasd",
			contentType: "text/plain",
			want: want{
				code: 400,
				location: "",
			},

		},
		{
			name: "test with empty {id} and contet-type: text/plain",
			shorturl: "",
			contentType: "text/plain",
			want: want{
				code: 400,
				location: "",
			},

		},

		{
			name: "test with correct shortUrl and content-type: '' ",
			shorturl: shortUrl,
			contentType: "",
			want: want{
				code: 400,
				location: "", 
			},
		},


	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Content-Type", test.contentType)
			req.SetPathValue("id", test.shorturl)

			wr := httptest.NewRecorder()
			
			app.RedirectToOriginalUrl(wr, req)
					
			res := wr.Result()
	
			assert.Equal(t, res.StatusCode, test.want.code)
			assert.Equal(t, res.Header.Get("Location"), test.want.location)

		})
	}
}


func TestReturnShortUrl(t *testing.T) {
	type want struct {
		code int
		contentType string
	}

	tests := []struct {
		name string
		method string
		body []byte
		contentType string
		want want
	}{
		{
			name: "test with correct request",
			method: http.MethodPost,
			body: []byte("https://www.google.com/"),
			contentType: "text/plain",
			want: want{
				code: 201,
				contentType: "text/plain",
			},
		},
		{
			name: "test with content-type: ''",
			method: http.MethodPost,
			body: []byte("https://hh.ru/"),
			contentType: "",
			want: want{
				code: 400,
				contentType: "",
			},
		},
		{
			name: "test with empty body",
			method: http.MethodPost,
			body: []byte(""),
			contentType: "text/plain",
			want: want{
				code: 400,
				contentType: "",
			},
		},
		{
			name: "test with request.Method NOT POST ",
			method: http.MethodGet,
			body: []byte(""),
			contentType: "text/plain",
			want: want{
				code: 400,
				contentType: "",
			},
		},

	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBuffer(test.body)
			
			req := httptest.NewRequest(test.method, "/", buf)
			req.Header.Set("Content-Type", test.contentType)
			
			wr := httptest.NewRecorder()
			
			app.ReturnShortUrl(wr, req)
			
			res := wr.Result()
			
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))	
		})
	}
}