package fliptext

import (
	"testing"
)

var teststring string = "ABCDEFGH"

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
	// for i := 0; i < len(teststring); i++ {
	// 	t.Logf("%v: %v\n", i, RotateTextNaiveLessMem(teststring, i))
	// }
	t.Log(RotateTextNaiveLessMem(teststring, 3))
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
