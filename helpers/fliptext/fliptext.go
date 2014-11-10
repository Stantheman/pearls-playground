package fliptext

import (
	"bytes"
	"github.com/Stantheman/pearls/helpers/gcd"
	// "math/big"
	// "os"
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
	if pointer == 0 {
		return t
	}
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

func RotateTextNaiveLessMem(t string, shift int) string {
	shift = shift % len(t)
	if shift == 0 {
		return t
	}
	return t[shift:] + t[:shift]
}

// for a different approach, we could define a subroutine to rotate X left one
// position, and call it I times, but that would be too time intensive
func RotateTextTimes(t string, times int) string {
	if times == 0 {
		return t
	}
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
	if shift == 0 {
		return t
	}
	perms := int(gcd.EuclidGCD(uint(shift), uint(length)))
	// make a modifiable copy of string
	bytes := make([]byte, length)
	copy(bytes, t)

	for perm := 0; perm < perms; perm++ {
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
func RotateDelicateBytes(bytes []byte, shift int) []byte {
	temp := bytes[0]
	length := len(bytes)
	shift = shift % length
	if shift == 0 {
		return bytes
	}
	perms := int(gcd.EuclidGCD(uint(shift), uint(length)))

	for perm := 0; perm < perms; perm++ {
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

	return bytes
}

func RotateReverseBytes(bytes []byte, shift int) []byte {
	shift = shift % len(bytes)
	if shift == 0 {
		return bytes
	}
	bytes = reverseBytes(bytes, 0, shift-1)
	bytes = reverseBytes(bytes, shift, len(bytes)-1)
	return reverseBytes(bytes, 0, len(bytes)-1)
}

func reverseBytes(bytes []byte, start, end int) []byte {
	for i, j := start, end; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return bytes
}
