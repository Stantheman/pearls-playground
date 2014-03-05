package io

import (
//	"fmt"
	"os"
)

// WriteFile checks if the filename exists, and if not, creates it
func WriteFile(filename string, data []byte) (bytesWritten int, err error) {

	bytesWritten = 0

	if _, err := os.Stat(filename); os.IsExist(err) {
		err = "refusing to write to " + filename + " since it exists"
		return bytesWritten, err
	}
	bytesWritten = 0
	return
}
