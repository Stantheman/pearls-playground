// Trying Programming Pearls's first example
package main

import (
	"bufio"
	"fmt"
	"github.com/Stantheman/pearls/helpers/bitmap"
	"github.com/Stantheman/pearls/helpers/random"
	"math"
	"os"
	"strconv"
)

func main() {

	const inputSize = 50

	input_file := "/tmp/wtf.txt"
	integers := random.GenerateUniqueRandomIntegers(inputSize)

	err := writeIntArray(integers, input_file, 1)
	if err != nil {
		fmt.Printf("Couldn't write out integers to file %v: %v", input_file, err)
		return
	}

	var thing bitmap.FirstBitmap
	thing.Sort(input_file, inputSize)
	fmt.Printf("thing is %v, %T\n", thing, thing)

	//originalProblem(inputSize)
	//firstProblem(inputSize)

}

const (
	input_file  = "/tmp/first_pearl_input.txt"
	output_file = "/tmp/first_pearl_output.txt"
)

func writeIntArray(integers []int, filename string, truncate int) (err error) {

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

/* The originalProblem is to take a list of unique integers in the range
1-27,000 from disk, sort them, and pump them out to disk using minimal memory.
The realization is that since the range is limited, known, and no metadata is attached,
a bitmap allows for us to read the input once and use a single bit to mark the
existance of an integer in the file.

This code creates the array of integers, writes then to a file delimited by newlines
and then procedes to recreate the solution listed above. */
func originalProblem(inputSize int) {

	integers := random.GenerateUniqueRandomIntegers(inputSize)

	err := writeIntArray(integers, input_file, 1)
	if err != nil {
		fmt.Printf("Couldn't write out integers to file %v: %v", input_file, err)
		return
	}
	/* this is related to question 2 -- what to do if your langauge doesn't support
	bitmaps -- using logical operations and otherwise */
	bitmap := make([]int8, inputSize)

	// now, pretending that we got this file externally, create the bitmap
	// and read it in
	infh, err := os.Open(input_file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer infh.Close()

	scanner := bufio.NewScanner(infh)
	for scanner.Scan() {
		// this is weird fake strconv -> int
		val, err := strconv.ParseInt(scanner.Text(), 10, 32)
		if err != nil {
			fmt.Printf("%v isn't a valid 32 bit integer: %v\n", val, err)
			return
		}
		bitmap[val] = 1
	}

	/* The original question had an input file and output file, but
	I needed an excuse to write the random number generation in go (vs opening stdin).*/
	// for i, v := range bitmap {
	// 	if v == 1 {
	// 		fmt.Println(i)
	// 	}
	// }
}

// the first problem is how to solve the issue if you have a limited amount of memory
func firstProblem(inputSize int) {

	integers := random.GenerateUniqueRandomIntegers(inputSize)

	// create the input file that will be read
	err := writeIntArray(integers, input_file, 1)
	if err != nil {
		fmt.Printf("Couldn't write integers to %v: %v\n", input_file, err)
		return
	}

	availableRam := inputSize / 5

	passes := int(math.Ceil(float64(inputSize) / float64(availableRam)))

	// open the input file for reading
	readFH, err := os.Open(input_file)
	if err != nil {
		fmt.Printf("Couldn't open %v for reading: %v\n", err)
		return
	}
	defer readFH.Close()

	// we're trading IO to keep memory low by reading multiple times
	bitmap := make([]int8, availableRam)
	for i := 0; i < passes; i++ {
		// get to the beginning of the file for each full pass
		_, err := readFH.Seek(0, 0)
		if err != nil {
			fmt.Printf("Unable to seek to 0 for %v: %v\n", input_file, err)
			return
		}

		scanner := bufio.NewScanner(readFH)
		for scanner.Scan() {
			// consume the number
			val, err := strconv.ParseInt(scanner.Text(), 10, 32)
			if err != nil {
				fmt.Printf("%v isn't a valid 32 bit integer: %v\n", val, err)
				return
			}

			// in the first pass, we only want to look at integers that are 0 < x < availableRam
			if int(val) >= (i*availableRam) && int(val) < (i*availableRam+availableRam) {
				bitmap[int(val)-(i*availableRam)] = 1
			}
		}

		//now that we've looped over the file, let's append the bitmap knowledge to our file
		// for index, value := range bitmap {
		// 	if value == 1 {
		// 		fmt.Println(index + (i * availableRam))
		// 	}
		// }
	}
}
