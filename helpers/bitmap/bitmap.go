// Package bitmap implements several bitmap-based sorts
//
// Each of these functions take an input file with newline-delimited integers < length.
// It iterates over the input, does its sort, and then dumps them to a file. Most also use
// an "avail" amount of space -- looping multiple times if the available space is less
// than the length of numbers.
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

// NaiveSort sorts using a slice of integers to store a single bit.
//
// NaiveSort is a helplessly weak attempt at mimicking the spirit of the first column of Programming Pearls.
// It was my first pass at learning Go and mostly kept around for prosterity.
//
// NaiveSort doesn't consider limited memory.
func NaiveSort(input_fn, output_fn string, length, avail int) (err error) {

	in, out, err := sortSetup(input_fn, output_fn, length, avail)
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
		if val < 0 {
			return fmt.Errorf("%v can't be less than 1\n", val)
		}
		if int(val) > length {
			return fmt.Errorf("%v can't be larger than %v\n", val, length)
		}
		if bits[val] == 1 {
			return fmt.Errorf("Duplicate input: we've already seen %v\n", val)
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

// LimitedSort is NaiveSort with memory constraints.
// This is the first real problem from the book and asks how we'd accomplish the task
// with less memory than the input size. The answer is to make multiple passes over the file,
// filling the bitmap each pass, then dumping to a file.
// This takes passes * input time, but does it with input/passes space
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

			if val < 0 {
				return fmt.Errorf("%v can't be less than 1\n", val)
			}
			if int(val) > length {
				return fmt.Errorf("%v can't be larger than %v\n", val, length)
			}

			// in the first pass, we only want to look at integers that are 0 < x < availableRam
			if int(val) >= min && int(val) < max {
				if bits[int(val)-min] == 1 {
					return fmt.Errorf("Duplicate input: we've already seen %v\n", val)
				}
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

// BitSort uses math/big's bitwise manipulation functions to sort an input file and write to an output file
// Problem #2 discusses the implementation of bitmaps. You can use a language's builtin
// support for bitmap operations, or implement your own with bitwise operations.
//
// It ends by asking how you'd implement these answers in COBOL and Pascal, which I'm skipping.
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
			if val < 0 {
				return fmt.Errorf("%v can't be less than 1\n", val)
			}
			if int(val) > length_b {
				return fmt.Errorf("%v can't be larger than %v\n", val, length_b)
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
				if bits[position].Bit(bit) == 1 {
					return fmt.Errorf("Duplicate input: we've already seen %v\n", val)
				}
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

// BitSortPrimative uses Go's builtin bitwise operations instead of math/big to sort
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
			if val < 0 {
				return fmt.Errorf("%v can't be less than 1\n", val)
			}
			if int(val) > length_b {
				return fmt.Errorf("%v can't be larger than %v\n", val, length_b)
			}

			// in the first pass, we only want to look at integers that are 0 < x < availableRam
			if int(val) >= min && int(val) < max {
				/* If we have 100 bits per pass and we see 125, we're on the second pass, so min = 1*100 = 100
				125-100 = 25. we're in the 25th bit slot, and 25/64 = 0, so the first available integer in the array*/
				position := (int(val) - min) >> 6
				/* If we have 10 bits per pass and see 68, it means we're on the 6th pass, which makes min = 6*10 = 60
				position is (68 - 60) / 64= integer 0. 68 - 0 - 60 = 8th bit */
				bit := int(val) - (position << 6) - min
				// fmt.Printf("pos is %v, bit is %v, value is %v\n", position, bit, val)
				if (bits[position] & (1 << uint(bit))) > 0 {
					return fmt.Errorf("Duplicate input: we've already seen %v\n", val)
				}
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

/* SortNonUnique sorts a list of integers from 1-N, where each value can occur M times

Problem 5 asks to solve the problem in the case where the programmer needs to sort
a list of 1-27000, but each value may happen up to 10 times. This requires 4 bits per
value instead of a single bit. This solution should allow for an arbitrary number of duplicates
up to 2^64*/

// this thing needs way more work and i'm gonna skip for now
// when i come back to this, remember,
// "When in doubt, use brute force" ~Ken Thompson
// func SortNonUnique(input_fn, output_fn string, length_b, avail_b, occur int) (err error) {

// 	// get set up
// 	if occur <= 0 {
// 		return fmt.Errorf("occur must be greater than 0: %v\n", occur)
// 	}
// 	in, out, err := sortSetup(input_fn, output_fn, length_b, avail_b)
// 	if err != nil {
// 		return err
// 	}
// 	defer in.Close()
// 	defer out.Close()
// 	writer := bufio.NewWriter(out)

// 	/* need to figure out how many ints we can do per run.
// 	if an integer can occur 10 times per pass, you need 4 bits to store that knowledge
// 	if you have 100 bits available per run, you can do 100/4 = 25 integers per pass.

// 	figure out the nearest upper power of two so we know how many integers
// 	we can stuff into a 64bit slot at a time. bit twiddling hacks has a cool trick
// 	to do this, but the easily accessible version is to take log base 2 of the number and add 1.
// 	log 2 69 = 6.xxx, =~ 7
// 	*/
// 	size_b := int(math.Log2(occur)) + 1
// 	if size_b > 64 {
// 		return fmt.Errorf("occur must be less than 2^64, size is %v\n", size_b)
// 	}
// 	valPerSlot := 64 / size
// 	// if we had 10 occurences, that's 4 bits per value, 16 values per 64bit int
// 	// if we have 128 bits available, we can do 32 values at a time (avail_b/size_b).
// 	valPerRun := avail_b / size_b
// 	slots := avail_b / 64
// 	// if we can do 32 values per run, we need length/valPerRun passes
// 	passes := length_b / valPerRun

// 	for i := 0; i < passes; i++ {

// 		// if we have 128 bits available, we need ceil(128/64) = 2 int64s to work with
// 		bits := make([]uint64, slots)
// 		for i := range bits {
// 			bits[i] = 0
// 		}

// 		scanner := bufio.NewScanner(in)
// 		min, max := i*valPerRun, i*valPerRun+valPerRun

// 		for scanner.Scan() {
// 			// consume the number
// 			val, err := strconv.ParseInt(scanner.Text(), 10, 32)
// 			if err != nil {
// 				return err
// 			}

// 			// in the first pass, we only want to look at integers that are 0 < x < valPerSlot
// 			if int(val) >= min && int(val) < max {
// 				slotNumber := (int(val) - min) >> 6
// 				// if the val is 3 and has 4 bits per val, we want to kick in at bit 12
// 				bitStart := val * size_b
// 				currentVal := bits[slotNumber] & ((1 << size_b - 1) << bitStart) >> bitStart

// 				if (currentVal == (1<<size_b-1)) {
// 					return fmt.Errorf("Too many %v values seen, bailing\n", val)
// 				}
// 				currentVal++

// 				bits[slotNumber] = bits[slotNumber] |

// 				/* we need to figure out what the current value is and increment*/

// 				/* If we have 100 bits per pass and we see 125, we're on the second pass, so min = 1*100 = 100
// 				125-100 = 25. we're in the 25th bit slot, and 25/64 = 0, so the first available integer in the array*/

// 				/* If we have 10 bits per pass and see 68, it means we're on the 6th pass, which makes min = 6*10 = 60
// 				position is (68 - 60) / 64= integer 0. 68 - 0 - 60 = 8th bit */
// 				bit := int(val) - (position << 6) - min
// 				// fmt.Printf("pos is %v, bit is %v, value is %v\n", position, bit, val)
// 				bits[position] = bits[position] | 1<<uint(bit)
// 			}
// 		}

// 		//now that we've looped over the file, let's append the bitmap knowledge to our file
// 		for i, v := range bits {
// 			for j := 0; j < 64; j++ {
// 				if calculated := v & (1 << uint(j)); calculated > 0 {
// 					/* If we have 20 bits to play with and this is the second loop, we're looking at 20-39
// 					20 bits = 1 64-bit integer to hold data. if bit 10 is set, we're at (0 * 64 + 10 + 20) = 30th bit*/
// 					value := i*64 + j + min
// 					_, err := fmt.Fprintf(writer, "%v\n", value)
// 					if err != nil {
// 						return err
// 					}
// 				}
// 			}
// 		}
// 		writer.Flush()

// 		// get to the beginning of the file for each full pass
// 		_, err := in.Seek(0, 0)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
