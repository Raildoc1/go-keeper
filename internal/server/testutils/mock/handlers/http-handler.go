package handlers

import (
	"io"
	"net/http"
)

var _ http.Handler = (*HTTPHandlerMock)(nil)

type HTTPHandlerMock struct{}

func NewHTTPHandlerMock() *HTTPHandlerMock {
	return &HTTPHandlerMock{}
}

func (H *HTTPHandlerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
