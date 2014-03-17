/* Tests for the bitmap functions.
Useful for fair total benchmarking:

go test -bench="OnePass|Go|Naive" -benchmem | column -t

Running something like the following two lines lets you dial in on a specific func:

./bitmap.test -test.bench="BenchmarkBitSortPrimative" -test.memprofile=mem.out -test.memprofilerate=1 -test.run="x"
go tool pprof --show_bytes --lines bitmap.test mem.out

*/
package bitmap

import (
	"bufio"
	"fmt"
	"github.com/Stantheman/pearls/helpers/random"
	"os"
	"sort"
	"strconv"
	"testing"
)

const (
	inputFile  = "/tmp/first_pearl_input.txt"
	outputFile = "/tmp/first_pearl_output.txt"
	inputSize  = 10000
	available  = inputSize / 10
)

func BenchmarkGoSort(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nums, err := readIntSlice(inputFile)
		if err != nil {
			return
		}
		sort.Ints(nums)
		writeIntSlice(nums, outputFile, 1)
	}
}

func TestNaiveSort(t *testing.T) {
	setup()

	err := NaiveSort(inputFile, outputFile, inputSize)
	if err != nil {
		t.Error(err)
	}

	compareInputAndOutput(t)
}

func BenchmarkNaiveSort(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NaiveSort(inputFile, outputFile, inputSize)
	}
}

func TestLimitedSort(t *testing.T) {
	setup()

	err := LimitedSort(inputFile, outputFile, inputSize, available)
	if err != nil {
		t.Error(err)
	}
	compareInputAndOutput(t)
}

func BenchmarkLimitedSort(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LimitedSort(inputFile, outputFile, inputSize, available)
	}
}

func BenchmarkLimitedSortOnePass(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LimitedSort(inputFile, outputFile, inputSize, inputSize)
	}
}

func TestBitSort(t *testing.T) {
	setup()

	err := BitSort(inputFile, outputFile, inputSize, available)
	if err != nil {
		t.Error(err)
	}
	compareInputAndOutput(t)
}

func BenchmarkBitSort(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BitSort(inputFile, outputFile, inputSize, available)
	}
}

func BenchmarkBitSortOnePass(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BitSort(inputFile, outputFile, inputSize, inputSize)
	}
}

func TestBitSortPrimative(t *testing.T) {
	setup()

	err := BitSortPrimative(inputFile, outputFile, inputSize, available)
	if err != nil {
		t.Error(err)
	}
	compareInputAndOutput(t)
}

func BenchmarkBitSortPrimative(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BitSortPrimative(inputFile, outputFile, inputSize, available)
	}
}

func BenchmarkBitSortPrimativeOnePass(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BitSortPrimative(inputFile, outputFile, inputSize, inputSize)
	}
}

func setup() {
	integers := random.GenerateUniqueRandomIntegers(inputSize)

	err := writeIntSlice(integers, inputFile, 1)
	if err != nil {
		fmt.Printf("Couldn't write out integers to file %v: %v", inputFile, err)
		return
	}
}

func compareInputAndOutput(t *testing.T) {
	// read the before and after files  to compare
	unsorted, err := readIntSlice(inputFile)
	if err != nil {
		t.Error(err)
	}

	sorted, _ := readIntSlice(outputFile)
	if err != nil {
		t.Error(err)
	}

	// check that the output is actually sorted
	if !sort.IntsAreSorted(sorted) {
		t.Error("The output isn't sorted")
	}

	// sort the input ints with a trusted sort to compare
	sort.Ints(unsorted)

	if len(unsorted) != len(sorted) {
		t.Errorf("The input length (%v) and output length (%v) aren't the same", len(unsorted), len(sorted))
	}

	for i := range sorted {
		if sorted[i] != unsorted[i] {
			t.Errorf("The value at %v isn't the same. Trusted: %v, Sort: %v", i, unsorted[i], sorted[i])
		}
	}
}

func writeIntSlice(integers []int, filename string, truncate int) (err error) {

	var writeFH *os.File

	if truncate != 0 {
		writeFH, err = os.Create(filename)
		if err != nil {
			return err
		}
	} else {
		_, err := os.Stat(filename)
		if err != nil {
			return err
		}
		writeFH, err = os.Open(filename)
		if err != nil {
			return err
		}

	}
	defer writeFH.Close()

	writer := bufio.NewWriter(writeFH)
	for _, v := range integers {
		_, err := fmt.Fprintf(writer, "%v\n", v)
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

func readIntSlice(input_fn string) (bits []int, err error) {
	in, err := os.Open(input_fn)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		// this is weird fake strconv -> int
		val, err := strconv.ParseInt(scanner.Text(), 10, 32)
		if err != nil {
			return nil, fmt.Errorf("%v isn't a valid 32 bit integer: %v\n", val, err)
		}
		bits = append(bits, int(val))
	}
	return bits, nil
}
