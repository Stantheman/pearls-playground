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
	inputSize  = 1000
	available  = inputSize / 10
)

type sorter func(string, string, int, int) error

// BenchmarkGoSort runs as a control to compare sort speeds
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

// this allows for meta-tests
var sorts = map[string]sorter{
	"NaiveSort":        NaiveSort,
	"BitSort":          BitSort,
	"LimitedSort":      LimitedSort,
	"BitSortPrimative": BitSortPrimative,
	"BitSortGo":        BitSortGo,
}

/* TestSort loops over the available sorts and runs test functions
with all the required scaffolding. This is kind of anti-testing
as far as Go is concerned...not sure if I'm proud of this or will end up reverting*/
func TestSort(t *testing.T) {
	for name, function := range sorts {
		if err := testSort(function)(t); err != nil {
			t.Errorf("%v: %v", name, err)
		}
		if err := testBrokenSort(function)(t); err != nil {
			t.Errorf("%v: %v", name, err)
		}
	}
}

func TestSortNonUnique(t *testing.T) {
	integers := random.GenerateLimitedRandomIntegers(inputSize, 10)

	err := writeIntSlice(integers, inputFile, 1)
	if err != nil {
		fmt.Printf("Couldn't write out integers to file %v: %v", inputFile, err)
		return
	}
	err = SortNonUnique(inputFile, outputFile, inputSize, available, 10)
	if err != nil {
		t.Error(err)
	}

	if err = compareInputAndOutput(); err != nil {
		t.Error(err)
	}
}

/* testSort returns a function that performs a basic test on the requested sort method
It generates an input file, sorts, then compares to Go sort to make sure it's sane*/
func testSort(f sorter) func(*testing.T) error {
	return func(t *testing.T) (err error) {
		setup()

		err = f(inputFile, outputFile, inputSize, available)
		if err != nil {
			return err
		}

		if err = compareInputAndOutput(); err != nil {
			return err
		}
		return nil
	}
}

/* testBrokenSort returns a function that generates a crappy input file
and makes sure that the sort barfs */
func testBrokenSort(f sorter) func(*testing.T) error {
	return func(t *testing.T) (err error) {
		breakingSetup()

		err = f(inputFile, outputFile, inputSize, available)
		if err == nil {
			return fmt.Errorf("Sort accepted bad input without returning error")
		}
		return nil
	}
}

/* TestBenchmarkSorts is probably something I actually feel guilty about. I wanted meta
benchmarks without copying and pasting code. The Go benchmark code doesn't like benchmarks
within benchmarks. I wanted to be able to still run "go test" and get individual
benchmark lines, as if each function was defined on its own. The workaround is to
name the function Test* so no benchmark is set up, then run through them*/
func TestBenchmarkSorts(t *testing.T) {
	for name, function := range sorts {
		// test sorts with default knowledge
		result := testing.Benchmark(benchmarkSort(function))
		t.Logf("%v: %v\t%v\n", name, result.String(), result.MemString())

		// test sorts with one pass
		result = testing.Benchmark(benchmarkSortOnePass(function))
		t.Logf("%v,one_pass: %v\t%v\n", name, result.String(), result.MemString())
	}
}

/* benchmarkSort generates a function that runs a benchmark on a given sort*/
func benchmarkSort(f sorter) func(*testing.B) {
	return func(b *testing.B) {
		setup()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f(inputFile, outputFile, inputSize, available)
		}
	}
}

/* benchmarkSortOnePass generates a function that benchmarks the given sort,
providing the sort with enough memory to sort in one pass*/
func benchmarkSortOnePass(f sorter) func(*testing.B) {
	return func(b *testing.B) {
		setup()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f(inputFile, outputFile, inputSize, inputSize)
		}
	}
}

// setup create a suitable input file for sort consumption
func setup() {
	integers := random.GenerateUniqueRandomIntegers(inputSize)

	err := writeIntSlice(integers, inputFile, 1)
	if err != nil {
		fmt.Printf("Couldn't write out integers to file %v: %v", inputFile, err)
		return
	}
}

/* breakingSetup creates artificially bad input
problem 6
- what if there's more than one integer in the input?
- what if the input is > N or < 0
- what should it do?
*/
func breakingSetup() {
	// 10-1 with two 7s
	integers := []int{10, 9, 8, 7, 7, 6, 5, 4, 3, 2, 1, -3}

	err := writeIntSlice(integers, inputFile, 1)
	if err != nil {
		fmt.Printf("Couldn't write out integers to file %v: %v", inputFile, err)
		return
	}
}

/* compareIinputAndOutput is run after a sort to make sure things actually got sorted*/
func compareInputAndOutput() (err error) {
	// read the before and after files  to compare
	unsorted, err := readIntSlice(inputFile)
	if err != nil {
		return err
	}

	sorted, err := readIntSlice(outputFile)
	if err != nil {
		return err
	}

	// check that the output is actually sorted
	if !sort.IntsAreSorted(sorted) {
		return fmt.Errorf("The output isn't sorted")
	}

	// sort the input ints with a trusted sort to compare
	sort.Ints(unsorted)

	if len(unsorted) != len(sorted) {
		return fmt.Errorf("The input length (%v) and output length (%v) aren't the same", len(unsorted), len(sorted))
	}

	for i := range sorted {
		if sorted[i] != unsorted[i] {
			return fmt.Errorf("The value at %v isn't the same. Trusted: %v, Sort: %v", i, unsorted[i], sorted[i])
		}
	}
	return nil
}

/* writeIntSlice is a weird way of writing integers out. Probably should have used more bufio love*/
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

/* readIntSlice consumes the file in the way we expect, integers followed by newlines*/
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
