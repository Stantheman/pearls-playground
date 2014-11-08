package fliptext

import (
	"bytes"
	// "fmt"
	"math/big"
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
// X[2I + 1] to X[I + 1]. if this hasn't moved every entry, try again
// with 2
func RotateDelicate(t string, shift int) string {
	temp := t[0]
	length := len(t)
	shift = shift % length
	var gcd big.Int
	gcd.GCD(nil, nil, big.NewInt(int64(length)), big.NewInt(int64(shift)))
	// make a modifiable copy of string
	bytes := make([]byte, length)
	copy(bytes, t)

	for perm := 0; perm < int(gcd.Int64()); perm++ {
		for i := shift + perm; ; i = (i + shift) % length {
			var send int = (length + i - shift) % length

			if i != perm {
				bytes[send] = t[i]
			} else {
				bytes[send] = temp
				break
			}

		}
	}

	return string(bytes)
}
func RotateDelicateBytes(bytes []byte, shift int) string {
	temp := bytes[0]
	length := len(bytes)
	shift = shift % length
	var gcd big.Int
	gcd.GCD(nil, nil, big.NewInt(int64(length)), big.NewInt(int64(shift)))

	for perm := 0; perm < int(gcd.Int64()); perm++ {
		for i := shift + perm; ; i = (i + shift) % length {
			var send int = (length + i - shift) % length

			if i != perm {
				bytes[send] = bytes[i]
			} else {
				bytes[send] = temp
				break
			}
		}
	}

	return string(bytes)
}
