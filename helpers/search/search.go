/*
Package search provides methods for finding information in a dataset.

search is introduced in column 2 of Programming Pearls
*/
package search

import (
	"bufio"
	"encoding/binary"
	//"fmt"
	"io"
	"os"
	"strconv"
)

const (
	ones_filename  string = "ones.bin"
	zeros_filename string = "zeros.bin"
)

// Missing takes a list of 32-bit numbers and the maximum int size in bits
// and the number of integers on the file and returns a missing int
func Missing(filename string, max, count int, mask, position uint32) (missing uint32, err error) {

	// open a bunch of filehandles
	err, in_fh, one_fh, zero_fh := setup(filename, mask)
	if err != nil {
		return 0, err
	}
	defer in_fh.Close()
	defer one_fh.Close()
	defer zero_fh.Close()

	// buffer writes
	bone := bufio.NewWriter(one_fh)
	bzero := bufio.NewWriter(zero_fh)

	var first_num, ones, zeros uint32 = 0, 0, 0

	// read the numbers in and count if the first bit is one or zero
	for err = binary.Read(in_fh, binary.BigEndian, &first_num); err == nil; err = binary.Read(in_fh, binary.BigEndian, &first_num) {
		if first_num&mask == mask {
			ones++
			binary.Write(bone, binary.BigEndian, first_num)
		} else {
			zeros++
			binary.Write(bzero, binary.BigEndian, first_num)
		}
	}
	// bail if it wasn't just EOF
	if err != io.EOF {
		return 0, err
	}
	// flush any pending writes
	bone.Flush()
	bzero.Flush()

	// if either side is empty, we now know a number that is missing
	var missingno uint32 = 0
	if zeros == 0 {
		// if there's no zeros, then a missing number is our mask with the current position toggled
		missingno = mask ^ 1<<(position)
		return missingno, nil
	} else if ones == 0 {
		// if there's no ones, our mask represents a missing integer
		missingno = mask
		return missingno, nil
	}

	// pick the next iteration
	var winner string
	var nextmask uint32
	// the mask is the right side of the tree
	if ones < zeros {
		winner = tempfile(mask, ones_filename)
		// set the bit in front of the current position to ON
		nextmask = mask | (1 << (position + 1))
	} else {
		winner = tempfile(mask, zeros_filename)
		/*
		* in order to go down the left branch, figure out our intermediate state
		* which is the current mask with the last position toggled. then, to that state,
		* toggle the 1 in front of the leader to ON
		 */
		nextmask = mask ^ (1 << position)
		nextmask = nextmask | (1 << (position + 1))
	}

	// call ourselves with the next batch
	missing, err = Missing(winner, max, count, nextmask, position+1)
	if err != nil {
		return 0, err
	}

	return missing, nil
}

func setup(filename string, mask uint32) (err error, in_fh, one_fh, zero_fh *os.File) {

	in_fh, err = os.Open(filename)
	if err != nil {
		return err, nil, nil, nil
	}

	one_fh, err = os.Create(tempfile(mask, ones_filename))
	if err != nil {
		return err, nil, nil, nil
	}

	zero_fh, err = os.Create(tempfile(mask, zeros_filename))
	if err != nil {
		return err, nil, nil, nil
	}
	return nil, in_fh, one_fh, zero_fh
}

func tempfile(mask uint32, suffix string) string {
	return strconv.Itoa(int(mask)) + suffix
}
