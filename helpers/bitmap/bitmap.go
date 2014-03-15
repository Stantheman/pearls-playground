/* implements a minimal bitmap structure with sort
for Programming Pearls - problem 1 */
package bitmap

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

/* The precise problem statement from Programming Pearls. X = 27000:
input: a file containing at most X integers in the range 1..x
 it is fatal if there are duplicated, and no other data is associated with the integer
output: a sorted list in increasing order of the input integers
constraints: at most, 1000 16bit words

This is a very naive answer to the final discussed solution, in which the programmer
found 27000 total free bits to store each integer's status. It does not actually
use bit math yet, so the savings aren't there yet*/
func NaiveSort(input_fn, output_fn string, length int) (err error) {

	if length < 0 {
		return fmt.Errorf("Length must be greater than 0: %v", length)
	}

	in, err := os.Open(input_fn)
	if err != nil {
		return err
	}
	defer in.Close()

	// same for output
	out, err := os.Create(output_fn)
	if err != nil {
		return err
	}
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
This takes passes * input time, but does it with input/passes space */
func LimitedSort(input_fn, output_fn string, length, avail int) (err error) {

	if length < 0 {
		return fmt.Errorf("Length must be greater than 0: %v", length)
	}

	if avail < 0 {
		return fmt.Errorf("Avail must be greater than 0: %v\n", avail)
	}

	// open the input file for reading
	in, err := os.Open(input_fn)
	if err != nil {
		return err
	}
	defer in.Close()

	// same for output
	out, err := os.Create(output_fn)
	if err != nil {
		return err
	}
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
