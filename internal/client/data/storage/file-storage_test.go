package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

type TestData struct {
	Data []byte
	Map  map[string]string
}

func TestRecover(t *testing.T) {
	t.Run("Test Read Non-Existent File", func(t *testing.T) {
		mp, err := readFromFile("./non-existent-path")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(mp))
	})

	t.Run("Test Decode Invalid JSON", func(t *testing.T) {
		r := strings.NewReader("} invalid json {")
		_, err := decodeJSON(r)
		assert.Error(t, err)
	})

	t.Run("Test Decode Valid JSON", func(t *testing.T) {
		r := strings.NewReader(
			`{
					"key-1" : "value-1",
					"key-2" : "value-2",
					"key-3" : "value-3"
				}`)
		_, err := decodeJSON(r)
		assert.NoError(t, err)
	})

	t.Run("Test Decode Valid JSON with wrong data type", func(t *testing.T) {
		r := strings.NewReader(`{ "number" : 125 }`)
		_, err := decodeJSON(r)
		assert.Error(t, err)
	})
}

func TestFileStorage(t *testing.T) {
	withEmptyStorage(t, func(t *testing.T, fs *FileStorage) {
		t.Run("Test Has Value True", func(t *testing.T) {
			err := fs.set("test-key", "test-value")
			require.NoError(t, err)
			assert.True(t, fs.Has("test-key"))
		})
		t.Run("Test Get Value Success", func(t *testing.T) {
			val, err := fs.get("test-key")
			require.NoError(t, err)
			assert.Equal(t, "test-value", val)
		})
		t.Run("Test Has Value False", func(t *testing.T) {
			err := fs.Reset("test-key")
			require.NoError(t, err)
			assert.False(t, fs.Has("test-key"))
		})
		t.Run("Test Get Value Fail", func(t *testing.T) {
			_, err := fs.get("test-key")
			assert.ErrorIs(t, err, ErrNotFound)
		})
	})
}

func TestJSONHelpers(t *testing.T) {
	type TestData struct {
		Str string            `json:"str"`
		I   int               `json:"i"`
		Map map[string]string `json:"map"`
	}

	type WrongTestData struct {
		Str int     `json:"str"`
		I   string  `json:"i"`
		Map float64 `json:"map"`
	}

	td := TestData{
		Str: "test-value",
		I:   123,
		Map: map[string]string{
			"test-key1": "test-value1",
			"test-key2": "test-value2",
			"test-key3": "test-value3",
		},
	}

	withEmptyStorage(t, func(t *testing.T, fs *FileStorage) {
		t.Run("Test Set Value Success", func(t *testing.T) {
			err := Set[TestData](fs, "test-key", td)
			require.NoError(t, err)
		})
		t.Run("Test Get Value Success", func(t *testing.T) {
			val, err := Get[TestData](fs, "test-key")
			require.NoError(t, err)
			assert.EqualValues(t, td, val)
		})
		t.Run("Test Get Value Fail", func(t *testing.T) {
			_, err := Get[WrongTestData](fs, "test-key")
			assert.Error(t, err)
		})
		t.Run("Test Get Value Fail", func(t *testing.T) {
			_, err := Get[TestData](fs, "wrong-key")
			assert.ErrorIs(t, err, ErrNotFound)
		})
	})
}

func withEmptyStorage(t *testing.T, f func(t *testing.T, fs *FileStorage)) {
	_ = os.Remove("./test.str")
	defer os.Remove("./test.str")

	fs, err := NewFileStorage("./test.str")
	require.NoError(t, err)

	f(t, fs)
}
