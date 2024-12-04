package middleware

import "net/http"

type memoryWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (w *memoryWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *memoryWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}
