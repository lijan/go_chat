package main

import (
	"io"
	"testing"
)

func TestSaveFile(t *testing.T) {
	conn := io.Reader
	fname := "testasdasdasdsadsdasdas123 1231 dhn dfas dfsahldfsakhsdflkdhasfklhasdlfhkahldfa dsf "
	contentPos := len(fname) + 2
	saveFile(conn, fname, contentPos)

}
