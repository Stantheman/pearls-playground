// Trying Programming Pearls's first example
package main

import (
	"fmt"
	"bufio"
	"strconv"
	"os"
	"github.com/Stantheman/pearls/helpers/random"
)

func main() {

	originalProblem()
	firstProblem()

}

/* The originalProblem is to take a list of unique integers in the range
1-27,000 from disk, sort them, and pump them out to disk using minimal memory.
The realization is that since the range is limited, known, and no metadata is attached,
a bitmap allows for us to read the input once and use a single bit to mark the
existance of an integer in the file.

This code creates the array of integers, writes then to a file delimited by newlines
and then procedes to recreate the solution listed above. */
func originalProblem() {

	const (
		input_file  = "/tmp/first_pearl_input.txt"
		output_file = "/tmp/first_pearl_output.txt"
		input_size  = 27000
	)

	integers := random.GenerateUniqueRandomIntegers(input_size)

	// write the random integers out
	writefh, err := os.Create(input_file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer writefh.Close()
	
	w := bufio.NewWriter(writefh)
	for _, v := range(integers) {
		fmt.Fprintf(w, "%v\n", v)
	}
	w.Flush()

	// now, pretending that we got this file externally, create the bitmap
	// and read it in


	/* this is related to question 2 -- what to do if your langauge doesn't support
	bitmaps -- using logical operations and otherwise */
	bitmap := make([]int8, input_size)

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
	
	for i, v := range bitmap {
		if v == 1 {
			fmt.Println(i)
		}
	}
}


func firstProblem() {
}
