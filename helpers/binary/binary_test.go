package binary

import (
	//"bytes"
	"encoding/binary"
	//"github.com/Stantheman/pearls/helpers/random"
	"bufio"
	//"io"
	"math/rand"
	"os"
	"testing"
	"time"
)

var filename = "wow.txt"
var create_size = 10000
var ints = generateRandomInt32s(create_size)

// benchmark make files
func BenchmarkMakingBinaryFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		makeBinaryFile(false)
	}
}

func BenchmarkMakingBinaryFileBuffered(b *testing.B) {
	for i := 0; i < b.N; i++ {
		makeBinaryFile(true)
	}
}

//benchmark read files
func BenchmarkReadingBinaryFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		readBinaryFile(false)
	}
}

func BenchmarkReadingBinaryFileBuffered(b *testing.B) {
	for i := 0; i < b.N; i++ {
		readBinaryFile(true)
	}
}

// test making files
func TestMakingBinaryFile(t *testing.T) {
	if err := makeBinaryFile(false); err != nil {
		t.Error(err)
	}
}

func TestMakingBinaryFileBuffered(t *testing.T) {
	if err := makeBinaryFile(true); err != nil {
		t.Error(err)
	}
}

// test reading files
func TestReadingBinaryFiles(t *testing.T) {
	_, err := readBinaryFile(false)
	if err != nil {
		t.Error(err)
	}
}

func TestReadingBinaryFilesBuffered(t *testing.T) {
	_, err := readBinaryFile(true)
	if err != nil {
		t.Error(err)
	}
}

// actual work
func makeBinaryFile(buf bool) error {
	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()

	if buf {
		bfh := bufio.NewWriter(fh)
		if err := binary.Write(bfh, binary.BigEndian, ints); err != nil {
			return err
		}
		bfh.Flush()
	} else {
		if err := binary.Write(fh, binary.BigEndian, ints); err != nil {
			return err
		}
	}

	return nil
}

// helper
func generateRandomInt32s(count int) (list []int32) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	list = make([]int32, count)

	for i := 0; i < count; i++ {
		list[i] = r.Int31()
	}
	return
}

func readBinaryFile(buf bool) ([]int32, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	readints := make([]int32, create_size)

	if buf {
		bfh := bufio.NewReader(fh)

		if err := binary.Read(bfh, binary.BigEndian, readints); err != nil {
			return nil, err
		}
	} else {
		if err := binary.Read(fh, binary.BigEndian, readints); err != nil {
			return nil, err
		}
	}
	return readints, nil
}
