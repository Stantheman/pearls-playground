package fliptext

import (
	// "fmt"
	"bytes"
	// "strings"
)

var buffer bytes.Buffer

/* Challenge B from Programming Pearls is to rotate a 1D array
of N characters left M positions. So the string ABCDEFGH with M=3
should become DEFGHABC*/

// FlipTextNaive takes the cue from the introduction,
// "Simple code uses an N-element intermediate vector to do the job
// in N steps"
func RotateTextNaive(t string, shift int) string {
	var pointer int = shift % len(t)
	// faster concatenation, same idea as intermediate += string(t[pointer])
	buffer.Truncate(0)

	for i := 0; i < len(t); i++ {
		if pointer == len(t) {
			pointer = 0
		}
		buffer.WriteByte(t[pointer])
		pointer++
	}
	return buffer.String()
}

// then the less memory intensive naive one, copies shift characters
// into a temporary buffer, shifts the rest of the array over, then appends
// cheating since this is less raw not shifting make 2 versions
func RotateTextNaiveLessMem(t string, shift int) string {
	intermediate := t[:shift]
	t = t[shift:]
	t += intermediate
	return t
}

// for a different approach, we could define a subroutine to rotate X left one
// position, and call it I times, but that would be too time intensive
func RotateTextTimes(t string, times int) string {
	for i := 0; i < times; i++ {
		t = RotateTextNaiveLessMem(t, 1)
	}
	return t
}

// one successful approach is just a delicate juggling act:
// move X[1] to temporary var t, then move X[I + 1] to X[1],
// X[2I + 1] to X[I + 1].
func RotateDelicate(t string, shift int) string {
	//temp := t[0]
	return "Wow"

}