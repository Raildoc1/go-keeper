package middleware

import (
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-keeper/internal/common/compression"
	"go-keeper/internal/server/testutils"
	"go-keeper/internal/server/testutils/mock/handlers"
	"go-keeper/pkg/logging"
	"net/http"
	"testing"
)

func TestResponseCompressor(t *testing.T) {
	logger := logging.NewNopLogger()
	const compressionLevel = gzip.BestSpeed
	responseCompressor := NewResponseCompressor(logger, compressionLevel)

	handlerSetup := testutils.HandlerSetup{
		Handler: responseCompressor.CreateHandler(handlers.NewHTTPHandlerMock()),
		Method:  http.MethodPost,
		URL:     "/api/user/test",
	}

	tests := []testutils.HandlerTestData{
		{
			TestName:     "Success",
			HandlerSetup: handlerSetup,
			Body:         "test_input",
			Headers: map[string]string{
				"Accept-Encoding": "gzip",
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   testutils.MustGzip("test_input", compressionLevel),
			CustomBodyChecker: func(t *testing.T, bytes []byte) {
				decompressed, err := compression.Decompress(bytes)
				require.NoError(t, err)
				assert.Equal(t, []byte("test_input"), decompressed)
			},
		},
		{
			TestName:       "No Compression",
			HandlerSetup:   handlerSetup,
			Body:           "test_input",
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   "test_input",
		},
	}

	testutils.PerformHTTPHandlerTests(t, tests)
}
