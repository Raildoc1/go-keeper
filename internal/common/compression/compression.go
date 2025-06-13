package compression

import (
	"bytes"
	"compress/gzip"
	"io"
)

func Compress(b []byte, compressionLevel int) ([]byte, error) {
	var compressedBodyBuffer bytes.Buffer
	gzipWriter, err := gzip.NewWriterLevel(&compressedBodyBuffer, compressionLevel)
	if err != nil {
		return nil, err
	}
	defer gzipWriter.Close()

	_, err = gzipWriter.Write(b)
	if err != nil {
		return nil, err
	}

	err = gzipWriter.Flush()
	if err != nil {
		return nil, err
	}

	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	return compressedBodyBuffer.Bytes(), nil
}

func Decompress(b []byte) ([]byte, error) {
	gzipReader, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	decompressed, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}
