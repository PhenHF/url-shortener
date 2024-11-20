package httpcompress

import (
	"compress/gzip"
	"io"
	"net/http"
)


type compressWriter struct {
	w http.ResponseWriter
	zw *gzip.Writer
}


func (cw *compressWriter) Header() http.Header {
	return cw.w.Header()
}

func (cw *compressWriter) Write(p []byte) (int, error) {
	return cw.zw.Write(p)
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		cw.w.Header().Set("Content-Encoding", "gzip")
	}
	
	cw.w.WriteHeader(statusCode)
}

func (cw *compressWriter) Close() error {
	return cw.zw.Close()
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w: w,
		zw: gzip.NewWriter(w),
	}
}

type compressReader struct {
	r io.ReadCloser
	zr *gzip.Reader
}


func (cr compressReader) Read(p []byte) (n int, err error) {
	return cr.zr.Read(p)
}

func (cr *compressReader) Close() error {
    if err := cr.r.Close(); err != nil {
        return err
    }
    return cr.zr.Close()
} 

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r: r,
		zr: zr,
	}, nil
}