package testutils

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/stretchr/testify/assert"
)

type HandlerSetup struct {
	Handler http.Handler
	Method  string
	URL     string
}

type HandlerTestData struct {
	TestName          string
	HandlerSetup      HandlerSetup
	Body              string
	PathParams        map[string]string
	Headers           map[string]string
	ExpectedStatus    int
	ExpectedBody      string
	CustomBodyChecker func(t *testing.T, body []byte)
}

func createResponseAndRequest(data *HandlerTestData) (w *httptest.ResponseRecorder, r *http.Request) {
	r = httptest.NewRequest(data.HandlerSetup.Method, data.HandlerSetup.URL, bytes.NewBufferString(data.Body))
	w = httptest.NewRecorder()
	if data.Headers != nil {
		for k, v := range data.Headers {
			r.Header.Set(k, v)
		}
	}
	if data.PathParams != nil {
		rctx := chi.NewRouteContext()
		for key, value := range data.PathParams {
			rctx.URLParams.Add(key, value)
		}
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	}
	return w, r
}

func PerformHTTPHandlerTests(t *testing.T, testDatas []HandlerTestData) {
	t.Helper()
	for _, testData := range testDatas {
		performHTTPHandlerTest(t, &testData)
	}
}

func performHTTPHandlerTest(t *testing.T, testData *HandlerTestData) {
	t.Helper()
	t.Run(testData.TestName, func(t *testing.T) {
		w, r := createResponseAndRequest(testData)

		testData.HandlerSetup.Handler.ServeHTTP(w, r)

		assert.Equal(t, testData.ExpectedStatus, w.Code, testData.TestName)

		if testData.CustomBodyChecker != nil {
			testData.CustomBodyChecker(t, w.Body.Bytes())
		} else {
			assert.Equal(t, testData.ExpectedBody, w.Body.String(), testData.TestName)
		}
	})
}
