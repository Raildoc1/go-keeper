package compression

import (
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestCompressDecompress(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
	}{
		{
			name:  "Simple String",
			bytes: []byte("Hello, World!"),
		},
		{
			name:  "Simple String",
			bytes: []byte("test_input"),
		},
		{
			name:  "Empty String",
			bytes: []byte(""),
		},
		{
			name:  "~1MB Random",
			bytes: randomBytes(1_000_000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := Compress(tt.bytes, gzip.BestCompression)
			require.NoError(t, err)
			decompressed, err := Decompress(compressed)
			require.NoError(t, err)
			assert.Equal(t, tt.bytes, decompressed)
		})
	}
}

func randomBytes(length int) []byte {
	result := make([]byte, length)
	rand.Read(result)
	return result
}
