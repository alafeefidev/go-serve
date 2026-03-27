package middleware

import (
	"bytes"
	"net/http"
)

type wrapper struct {
	http.ResponseWriter
	buf        bytes.Buffer
}

func newResponseWriterWrapper(w http.ResponseWriter) *wrapper {
	return &wrapper{ResponseWriter: w}
}

func (w *wrapper) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}
