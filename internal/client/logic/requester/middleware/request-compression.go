package middleware

import (
	"compress/gzip"
	"github.com/go-resty/resty/v2"
	"go-keeper/internal/common/compression"
)

var _ resty.RequestMiddleware = (&CompressionMiddleware{}).Execute

type CompressionMiddleware struct {
}

func NewCompressionMiddleware() *CompressionMiddleware {
	return &CompressionMiddleware{}
}

func (m *CompressionMiddleware) Execute(_ *resty.Client, r *resty.Request) error {
	bodyBytes, ok := r.Body.([]byte)
	if !ok {
		return nil
	}

	const minBodyLengthToCompress = 512
	if len(bodyBytes) < minBodyLengthToCompress {
		return nil
	}

	compressedBody, err := compression.Compress(bodyBytes, gzip.BestSpeed)
	if err != nil {
		return err
	}

	r.SetHeader("Content-Encoding", "gzip")
	r.SetBody(compressedBody)

	return nil
}
