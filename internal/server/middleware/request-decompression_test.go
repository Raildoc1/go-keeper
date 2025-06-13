package middleware

import (
	"compress/gzip"
	"go-keeper/internal/server/testutils"
	"go-keeper/internal/server/testutils/mock/handlers"
	"go-keeper/pkg/logging"
	"net/http"
	"testing"
)

func TestRequestDecompressor(t *testing.T) {
	logger := logging.NewNopLogger()
	requestDecompressor := NewRequestDecompressor(logger)

	const compressionLevel = gzip.BestCompression

	handlerSetup := testutils.HandlerSetup{
		Handler: requestDecompressor.CreateHandler(handlers.NewHTTPHandlerMock()),
		Method:  http.MethodPost,
		URL:     "/api/user/test",
	}

	tests := []testutils.HandlerTestData{
		{
			TestName:     "Success",
			HandlerSetup: handlerSetup,
			Body:         testutils.MustGzip("test_input", compressionLevel),
			Headers: map[string]string{
				"Content-Encoding": "gzip",
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   "test_input",
		},
		{
			TestName:       "No Header",
			HandlerSetup:   handlerSetup,
			Body:           testutils.MustGzip("test_input", compressionLevel),
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   testutils.MustGzip("test_input", compressionLevel),
		},
		{
			TestName:     "Wrong Data Format",
			HandlerSetup: handlerSetup,
			Body:         "not gzip data",
			Headers: map[string]string{
				"Content-Encoding": "gzip",
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
		{
			TestName:       "Non-Compressed Data",
			HandlerSetup:   handlerSetup,
			Body:           "non-compressed input",
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   "non-compressed input",
		},
	}

	testutils.PerformHTTPHandlerTests(t, tests)
}
