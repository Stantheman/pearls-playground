/* implements a minimal bitmap structure with sort
for Programming Pearls - problem 1 */
package bitmap

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Bitmap interface {
	Sort(filename string, length int)
}

type FirstBitmap []int8

/* The precise problem statement from Programming Pearls. X = 27000:
input: a file containing at most X integers in the range 1..x
 it is fatal if there are duplicated, and no other data is associated with the integer
output: a sorted list in increasing order of the input integers
constraints: at most, 1000 16bit words

This is a very naive answer to the final discussed solution, in which the programmer
found 27000 total free bits to store each integer's status. It does not actually
use bit math yet, so the savings aren't there yet*/
func (b FirstBitmap) Sort(filename string, length int) error {

	fmt.Printf("bitmap is %v, type %T\n", b, b)
	if length < 0 {
		return fmt.Errorf("Length must be greater than 0: %v", length)
	}

	in, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer in.Close()

	b = make([]int8, length)
	fmt.Printf("bitmap is %v, type %T\n", b, b)
	//fmt.Printf("*bitmap is %v, type %T\n", *b, *b)
	//fmt.Printf("*bitmap is %v, type %T\n", &b, &b)

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		// this is weird fake strconv -> int
		val, err := strconv.ParseInt(scanner.Text(), 10, 32)
		if err != nil {
			return fmt.Errorf("%v isn't a valid 32 bit integer: %v\n", val, err)
		}
		b[val] = 1
	}
	fmt.Printf("bitmap is %v, type %T\n", b, b)
	return nil
}
