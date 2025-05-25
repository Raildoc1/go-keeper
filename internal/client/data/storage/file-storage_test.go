package storage

import "testing"

type TestData struct {
	Data []byte
	Map  map[string]string
}

func TestFileStorage(t *testing.T) {
	fs, err := NewFileStorage("./test.str")
	if err != nil {
		t.Fatal(err)
	}
	td := TestData{
		Data: []byte("hello world"),
		Map: map[string]string{
			"test-key-1": "test-value-1",
			"test-key-2": "test-value-2",
			"test-key-3": "test-value-3",
		},
	}
	err = Set(fs, "test-key", td)
	if err != nil {
		t.Fatal(err)
	}
	err = fs.Save("./test.str")
	if err != nil {
		t.Fatal(err)
	}
}
