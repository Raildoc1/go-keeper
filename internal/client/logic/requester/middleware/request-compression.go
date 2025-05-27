package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/go-resty/resty/v2"
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

	var compressedBodyBuffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedBodyBuffer)
	defer gzipWriter.Close()

	_, err := gzipWriter.Write(bodyBytes)
	if err != nil {
		return err
	}

	err = gzipWriter.Flush()
	if err != nil {
		return err
	}

	compressedB := compressedBodyBuffer.Bytes()

	r.SetHeader("Content-Encoding", "gzip")
	r.SetBody(compressedB)

	return nil
}
