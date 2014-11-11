package fliptext

import (
	"testing"
)

var teststring string = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccccccccccccccccccccccccccddddddddddddddddddddddddddddddddddddddddddddddddeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeffffffffffffffffffffffffffffffffffffffffffffffffgggggggggggggggggggggggggggggggggggggggggggggggghhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiijjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkllllllllllllllllllllllllllllllllllllllllllllllllmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnooooooooooooooooooooooooooooooooooooooooooooooooppppppppppppppppppppppppppppppppppppppppppppppppqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrssssssssssssssssssssssssssssssssssssssssssssssssttttttttttttttttttttttttttttttttttttttttttttttttuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"

func TestRotateTextNaive(t *testing.T) {
	for i := 0; i < len(teststring)*2; i++ {
		t.Logf("%v: %v\n", i, RotateTextNaive(teststring, i))
	}
}
func BenchmarkRotateTextNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RotateTextNaive(teststring, 3)
	}
}

func TestRotateTextNaiveLessMem(t *testing.T) {
	for i := 0; i < len(teststring)*2; i++ {
		t.Logf("%v: %v\n", i, RotateTextNaiveLessMem(teststring, i))
	}
}
func BenchmarkRotateTextNaiveLessMem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RotateTextNaiveLessMem(teststring, 3)
	}
}

func TestRotateTextTimes(t *testing.T) {
	for i := 0; i < len(teststring)*2; i++ {
		t.Logf("%v: %v\n", i, RotateTextTimes(teststring, i))
	}
}
func BenchmarkRotateTextTimes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RotateTextTimes(teststring, 3)
	}
}

func TestRotateDelicate(t *testing.T) {
	for i := 0; i < len(teststring)*2; i++ {
		t.Log(RotateDelicate(teststring, i))
	}
}
func BenchmarkRotateDelicate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RotateDelicate(teststring, 3)
	}
}

func TestRotateDelicateBytes(t *testing.T) {
	for i := 0; i < len(teststring)*2; i++ {
		t.Log(string(RotateDelicateBytes([]byte(teststring), i)))
	}
}
func BenchmarkRotateDelicateBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RotateDelicateBytes([]byte(teststring), 3)
	}
}

func TestRotateReverseBytes(t *testing.T) {
	for i := 0; i < len(teststring)*2; i++ {
		t.Log(string(RotateReverseBytes([]byte(teststring), i)))
	}
}
func BenchmarkRotateReverseBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RotateReverseBytes([]byte(teststring), 3)
	}
}
