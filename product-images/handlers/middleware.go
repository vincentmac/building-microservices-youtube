package handlers

import (
	"compress/gzip"
	"net/http"
)

// GZipResponseMiddleware detects if teh client can handle
// zipped content and if so returns the response in GZipped format
func GZipResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// write the file
		next.ServeHTTP(rw, r)
	})
}

type WrappedResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedResponseWriter {

	// wrap the default writer in a gzip writer
	gw := gzip.NewWriter(rw)

	return &WrappedResponseWriter{rw, gw}
}

func (wr *WrappedResponseWriter) Header() http.Header {
	return wr.rw.Header()
}

func (wr *WrappedResponseWriter) Write(d []byte) (int, error) {
	return wr.gw.Write(d)
}

func (wr *WrappedResponseWriter) WriteHeader(statusCode int) {
	wr.rw.WriteHeader(statusCode)
}
