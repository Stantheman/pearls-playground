package random

import (
	"math"
	"testing"
)

const (
	ArraySize = 50
)

// TestGenerateIncreasingRandomIntegers checks size, uniqueness, and increasingness
func TestGenerateIncreasingRandomIntegers(t *testing.T) {
	numbers := GenerateIncreasingRandomIntegers(ArraySize)
	length := len(numbers)

	if length != ArraySize {
		t.Errorf("Expected %v integers, got %v\n", ArraySize, length)
	}

	t.Logf("The %v numbers are: %v\n", length, numbers)

	previous := -1
	for _, v := range numbers {
		if v <= previous {
			t.Errorf("Integers aren't in increasing order. Saw %v, previous: %v\n", v, previous)
		}
		if v < 0 {
			t.Errorf("Integers must be positive. Saw %v\n", v)
		}
		previous = v
	}
}

// TestGenerateUniqueRandomIntegers checks size and uniqueness
func TestGenerateUniqueRandomIntegers(t *testing.T) {
	numbers := GenerateUniqueRandomIntegers(ArraySize)
	length := len(numbers)

	if length != ArraySize {
		t.Errorf("Expected %v integers, got %v\n", ArraySize, length)
	}

	t.Logf("The %v numbers are: %v\n", length, numbers)

	// terrible way of checking uniqueness, should sort and loop
	// wanna keep forcing myself to get used to typing this first
	for i, vi := range numbers {
		for j, vj := range numbers {
			if i == j {
				continue
			}
			if vi == vj {
				t.Errorf("Index %v and %v contain the same value %v\n", i, j, vi)
			}
		}
	}
}

func TestDistributionRandom(t *testing.T) {
	var numbers [ArraySize]int
	runs := 100000

	for i := 0; i < runs; i++ {
		list := GenerateUniqueRandomIntegers(ArraySize)
		for j := 0; j < ArraySize; j++ {
			numbers[j] += list[j]
		}
	}

	// average is the (summation of (ArraySize -1) * runs)/ArraySize
	average := ((ArraySize - 1) * ArraySize) / 2 * runs / ArraySize

	// stddev
	var stddevArray [ArraySize]float64
	for i, v := range numbers {
		stddevArray[i] = math.Pow(float64(v-average), 2.0)
	}

	total := 0.0
	for _, v := range stddevArray {
		total += v
	}
	stddev := math.Sqrt(total / ArraySize)
	t.Logf("std deviation is %v\n", stddev)

	// figure out how equally random we are
	// http://en.wikipedia.org/wiki/68-95-99.7_rule
	var distribution [3]int
	for _, v := range numbers {
		difference := int(math.Abs(float64(average-v)) / stddev)

		if difference < len(distribution) {
			distribution[difference]++
		} else {
			t.Errorf("%v is more than 3 standard devations away\n")
		}
	}

	// print the pretty knowledge
	for i := range distribution {
		t.Logf("%v%% fall within %v standard deviations\n", float64(distribution[i])/float64(ArraySize)*100, i+1)
	}
}

func BenchmarkGenerateIncreasingRandomIntegers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateIncreasingRandomIntegers(ArraySize)
	}
}

func BenchmarkGenerateUniqueRandomIntegers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateUniqueRandomIntegers(ArraySize)
	}
}
