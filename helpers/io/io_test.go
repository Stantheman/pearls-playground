package io

import (
	"testing"
)

// TestWriteFile checks to see if we could write to a file
func TestWriteFile(t *testing.T) {
	byte := []byte{'a', 'b', 'c'}
	WriteFile("/tmp/lol", byte)
}
