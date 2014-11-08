package gcd

import (
	"testing"
)

var tests = [][]uint{
	{10, 5, 5},
	{20, 10, 10},
	{0, 5, 5},
	{5, 0, 5},
	{1, 50, 1},
	{50, 1, 1},
	{100, 2, 2},
	{1, 1, 1},
}

func TestEuclidGCD(t *testing.T) {
	for i, set := range tests {
		if res := EuclidGCD(set[0], set[1]); res != set[2] {
			t.Errorf("%v: res is %v, expected %v\n", i, res, set[2])
		}
	}
}
func BenchmarkEuclidGCD(b *testing.B) {
	for j := 0; j < b.N; j++ {
		for i, set := range tests {
			if res := EuclidGCD(set[0], set[1]); res != set[2] {
				b.Errorf("%v: res is %v, expected %v\n", i, res, set[2])
			}
		}
	}
}

func TestBinaryGCD(t *testing.T) {
	for i, set := range tests {
		if res := BinaryGCD(set[0], set[1]); res != set[2] {
			t.Errorf("%v: res is %v, expected %v\n", i, res, set[2])
		}
	}
}
func BenchmarkBinaryGCD(b *testing.B) {
	for j := 0; j < b.N; j++ {
		for i, set := range tests {
			if res := BinaryGCD(set[0], set[1]); res != set[2] {
				b.Errorf("%v: res is %v, expected %v\n", i, res, set[2])
			}
		}
	}
}

func TestIterativeBinaryGCD(t *testing.T) {
	for i, set := range tests {
		if res := IterativeBinaryGCD(set[0], set[1]); res != set[2] {
			t.Errorf("%v: res is %v, expected %v\n", i, res, set[2])
		}
	}
}
func BenchmarkIterativeBinaryGCD(b *testing.B) {
	for j := 0; j < b.N; j++ {
		for i, set := range tests {
			if res := IterativeBinaryGCD(set[0], set[1]); res != set[2] {
				b.Errorf("%v: res is %v, expected %v\n", i, res, set[2])
			}
		}
	}
}
