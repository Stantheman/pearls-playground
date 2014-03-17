/* implements a minimal bitmap structure with sort
for Programming Pearls - problem 1 */
package bitmap

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
)

/* sortSetup is a helper sub to do the grunt work before the sort begins.
Simple error checking, returns two filehandles or possibly an error. Caller
must check error and is responsible for calling close or deferring */
func sortSetup(input_fn, output_fn string, length, avail int) (in, out *os.File, err error) {
	if length < 0 {
		return nil, nil, fmt.Errorf("Length must be greater than 0: %v", length)
	}

	if avail < 0 {
		return nil, nil, fmt.Errorf("Avail must be greater than 0: %v\n", avail)
	}

	// open the input file for reading
	in, err = os.Open(input_fn)
	if err != nil {
		return nil, nil, err
	}

	// same for output
	out, err = os.Create(output_fn)
	if err != nil {
		return nil, nil, err
	}

	return
}

/* The precise problem statement from Programming Pearls. X = 27000:
input: a file containing at most X integers in the range 1..x
 it is fatal if there are duplicated, and no other data is associated with the integer
output: a sorted list in increasing order of the input integers
constraints: at most, 1000 16bit words

This is a very naive answer to the final discussed solution, in which the programmer
found 27000 total free bits to store each integer's status. It does not actually
use bit math yet, so the savings aren't there yet.*/
func NaiveSort(input_fn, output_fn string, length int) (err error) {

	// NaiveSort is special since it's the only function without the need for available ram
	in, out, err := sortSetup(input_fn, output_fn, length, 1)
	if err != nil {
		return err
	}
	defer in.Close()
	defer out.Close()

	bits := make([]int, length)

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		// this is weird fake strconv -> int
		val, err := strconv.ParseInt(scanner.Text(), 10, 32)
		if err != nil {
			return fmt.Errorf("%v isn't a valid 32 bit integer: %v\n", val, err)
		}
		bits[val] = 1
	}

	writer := bufio.NewWriter(out)
	for i, v := range bits {
		if v == 1 {
			_, err := fmt.Fprintf(writer, "%v\n", i)
			if err != nil {
				return err
			}
		}
	}
	writer.Flush()
	return nil
}

/* LimitedSort takes a filename and a length of memory available.
This is the first real problem from the book and asks how we'd accomplish the task
with less memory than the input size. The answer is to make multiple passes over the file,
filling the bitmap each pass, then dumping to a file.
This takes passes * input time, but does it with input/passes space

In reality this can be NaiveSort where length = avail, but I wrote NaiveSort first*/
func LimitedSort(input_fn, output_fn string, length, avail int) (err error) {

	in, out, err := sortSetup(input_fn, output_fn, length, avail)
	if err != nil {
		return err
	}
	defer in.Close()
	defer out.Close()

	writer := bufio.NewWriter(out)
	passes := int(math.Ceil(float64(length) / float64(avail)))

	for i := 0; i < passes; i++ {

		bits := make([]int, avail)
		scanner := bufio.NewScanner(in)
		min, max := i*avail, i*avail+avail

		for scanner.Scan() {
			// consume the number
			val, err := strconv.ParseInt(scanner.Text(), 10, 32)
			if err != nil {
				return err
			}

			// in the first pass, we only want to look at integers that are 0 < x < availableRam
			if int(val) >= min && int(val) < max {
				bits[int(val)-min] = 1
			}
		}

		//now that we've looped over the file, let's append the bitmap knowledge to our file
		for i, v := range bits {
			if v == 1 {
				_, err := fmt.Fprintf(writer, "%v\n", i+min)
				if err != nil {
					return err
				}
			}
		}
		writer.Flush()

		// get to the beginning of the file for each full pass
		_, err := in.Seek(0, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

/* Problem #2 discusses the implementation of bitmaps. You can use a language's builtin
support for bitmap operations, or implement your own with bitwise operations.

It ends by asking how you'd implement these answers in COBOL and Pascal, which I'm skipping.

The goal here is to use math.big's bitmap support to answer problem 1. */
func BitSort(input_fn, output_fn string, length_b, avail_b int) (err error) {

	in, out, err := sortSetup(input_fn, output_fn, length_b, avail_b)
	if err != nil {
		return err
	}
	defer in.Close()
	defer out.Close()

	writer := bufio.NewWriter(out)
	passes := int(math.Ceil(float64(length_b) / float64(avail_b)))
	// fmt.Printf("passes is %v\n", passes)

	for i := 0; i < passes; i++ {

		// if we have 100 bits available, we need ceil(100/64) = ~2 int64s to work with
		bits := make([]*big.Int, int(math.Ceil(float64(avail_b)/64.0)))
		for i := range bits {
			bits[i] = big.NewInt(0)
		}

		// fmt.Printf("We made %v bits!\n", len(bits))

		scanner := bufio.NewScanner(in)
		min, max := i*avail_b, i*avail_b+avail_b

		// fmt.Printf("min: %v max: %v\n", min, max)
		// fmt.Printf("on pass %v\n", i)

		for scanner.Scan() {
			// consume the number
			val, err := strconv.ParseInt(scanner.Text(), 10, 32)
			if err != nil {
				return err
			}

			// in the first pass, we only want to look at integers that are 0 < x < availableRam
			if int(val) >= min && int(val) < max {
				/* If we have 100 bits per pass and we see 125, we're on the second pass, so min = 1*100 = 100
				125-100 = 25. we're in the 25th bit slot, and 25/64 = 0, so the first available integer in the array*/
				position := (int(val) - min) / 64
				/* If we have 10 bits per pass and see 68, it means we're on the 6th pass, which makes min = 6*10 = 60
				position is (68 - 60) / 64= integer 0. 68 - 0 - 60 = 8th bit */
				bit := int(val) - (64 * position) - min
				//fmt.Println(position, bit, val)
				bits[position].SetBit(bits[position], bit, 1)
			}
		}

		//now that we've looped over the file, let's append the bitmap knowledge to our file
		for i, v := range bits {
			for j := 0; j < 64; j++ {
				if v.Bit(j) == 1 {
					/* If we have 20 bits to play with and this is the second loop, we're looking at 20-39
					20 bits = 1 64-bit integer to hold data. if bit 10 is set, we're at (0 * 64 + 10 + 20) = 30th bit*/
					value := i*64 + j + min
					_, err := fmt.Fprintf(writer, "%v\n", value)
					if err != nil {
						return err
					}
				}
			}
		}
		writer.Flush()

		// get to the beginning of the file for each full pass
		_, err := in.Seek(0, 0)
		if err != nil {
			return err
		}
	}
	return nil

}

/* The second half of the second question asks about bitwise operations. This function
should provide the same functionality as BitSort without using the convenience functions
from math/big */
func BitSortPrimative(input_fn, output_fn string, length_b, avail_b int) (err error) {

	in, out, err := sortSetup(input_fn, output_fn, length_b, avail_b)
	if err != nil {
		return err
	}
	defer in.Close()
	defer out.Close()
	writer := bufio.NewWriter(out)

	passes := int(math.Ceil(float64(length_b) / float64(avail_b)))
	for i := 0; i < passes; i++ {

		// if we have 100 bits available, we need ceil(100/64) = ~2 int64s to work with
		//	bits := make([]*big.Int, int(math.Ceil(float64(avail_b)/64.0)))
		bits := make([]uint64, int(math.Ceil(float64(avail_b)/64.0)))
		for i := range bits {
			bits[i] = 0
		}

		scanner := bufio.NewScanner(in)
		min, max := i*avail_b, i*avail_b+avail_b

		for scanner.Scan() {
			// consume the number
			val, err := strconv.ParseInt(scanner.Text(), 10, 32)
			if err != nil {
				return err
			}

			// in the first pass, we only want to look at integers that are 0 < x < availableRam
			if int(val) >= min && int(val) < max {
				/* If we have 100 bits per pass and we see 125, we're on the second pass, so min = 1*100 = 100
				125-100 = 25. we're in the 25th bit slot, and 25/64 = 0, so the first available integer in the array*/
				position := (int(val) - min) / 64
				/* If we have 10 bits per pass and see 68, it means we're on the 6th pass, which makes min = 6*10 = 60
				position is (68 - 60) / 64= integer 0. 68 - 0 - 60 = 8th bit */
				bit := int(val) - (64 * position) - min
				// fmt.Printf("pos is %v, bit is %v, value is %v\n", position, bit, val)
				bits[position] = bits[position] | 1<<uint(bit)
			}
		}

		//now that we've looped over the file, let's append the bitmap knowledge to our file
		for i, v := range bits {
			for j := 0; j < 64; j++ {
				if calculated := v & (1 << uint(j)); calculated > 0 {
					/* If we have 20 bits to play with and this is the second loop, we're looking at 20-39
					20 bits = 1 64-bit integer to hold data. if bit 10 is set, we're at (0 * 64 + 10 + 20) = 30th bit*/
					value := i*64 + j + min
					_, err := fmt.Fprintf(writer, "%v\n", value)
					if err != nil {
						return err
					}
				}
			}
		}
		writer.Flush()

		// get to the beginning of the file for each full pass
		_, err := in.Seek(0, 0)
		if err != nil {
			return err
		}
	}
	return nil
}
