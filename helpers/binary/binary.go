package binary

import (
	"encoding/binary"
	"os"
)

func MakeBinaryFile(filename string, ints []uint32) error {
	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()

	if err := binary.Write(fh, binary.BigEndian, ints); err != nil {
		return err
	}

	return nil
}

func ReadBinaryFile(filename string) ([]uint32, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	readints := make([]uint32, 1)

	if err := binary.Read(fh, binary.BigEndian, readints); err != nil {
		return nil, err
	}

	return readints, nil
}
