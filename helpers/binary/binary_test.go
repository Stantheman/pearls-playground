package binary

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

var filename = "wow.txt"
var create_size = 10000
var ints = generateRandomUint32s(create_size)

// benchmark make files
func BenchmarkMakingBinaryFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MakeBinaryFile(filename, ints)
	}
}

// control benchmark
func BenchmarkMakingTextFileBuffered(b *testing.B) {
	for i := 0; i < b.N; i++ {
		writeIntSlice()
	}
}

//benchmark read files
func BenchmarkReadingBinaryFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadBinaryFile(filename)
	}
}

// test making files
func TestMakingBinaryFile(t *testing.T) {
	if err := MakeBinaryFile(filename, ints); err != nil {
		t.Error(err)
	}
}

// test reading files
func TestReadingBinaryFiles(t *testing.T) {
	_, err := ReadBinaryFile(filename)
	if err != nil {
		t.Error(err)
	}
}

// example control version
func writeIntSlice() (err error) {

	writeFH, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer writeFH.Close()

	writer := bufio.NewWriter(writeFH)
	for _, v := range ints {
		_, err := writer.WriteString(strconv.Itoa(int(v)) + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

// helper
func generateRandomUint32s(count int) (list []uint32) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	list = make([]uint32, count)

	for i := 0; i < count; i++ {
		list[i] = r.Uint32()
	}
	return
}
