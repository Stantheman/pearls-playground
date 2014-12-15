package random

import (
	"math"
	"sort"
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

func TestGenerateRandomIntegers(t *testing.T) {
	size := 5
	numbers, err := GenerateRandomIntegers(size, 32)
	if err != nil {
		t.Error(err)
	}
	if len(numbers) != size {
		t.Errorf("not the same size. result: %v, expected %v\n", len(numbers), size)
	}

	perc := make(map[uint32]int)
	for i := 0; i < 1000; i++ {

		numbers, err = GenerateRandomIntegers(1, 3)
		perc[numbers[0]] += 1
	}
	t.Logf("%#v\n", perc)
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

// TestGenerateUniqueRandomIntegers checks size and uniqueness
func TestGenerateLimitedRandomIntegers(t *testing.T) {
	limit := 10
	numbers := GenerateLimitedRandomIntegers(ArraySize, limit)
	length := len(numbers)

	// we don't know how many elements we're asking for, but we know it can't be more than ArraySize*limit
	if length > ArraySize*limit {
		t.Errorf("The length is larger than the possible outcome size: %v length, theoretical max of %v\n", length, ArraySize*limit)
	}

	sort.Ints(numbers)

	last, count := 0, 0
	for _, v := range numbers {
		if v == last {
			count++
		} else {
			count = 1
			last = v
		}
		if count > limit {
			t.Errorf("More than %v occurences of %v\n", limit, last)
		}
		if v > ArraySize {
			t.Errorf("%v is larger than the max of %v\n", v, ArraySize)
		}
	}
}

func TestDistributionRandom(t *testing.T) {
	var numbers [ArraySize]int
	runs := 100000

	for i := 0; i < runs; i++ {
		list, _ := GenerateRandomIntegers(ArraySize, 32)
		//t.Logf("list: %v, numbers %v\n", len(list), len(numbers))
		for j := 0; j < ArraySize; j++ {
			numbers[j] += int(list[j])
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
			t.Errorf("%v is more than 3 standard devations away\n", difference)
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

func BenchmarkGenerateLimitedRandomIntegers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateLimitedRandomIntegers(ArraySize, 10)
	}
}
