// Package random implements specific random number generator functions
package random

import (
	"errors"
	"math/rand"
	"sort"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateIncreasingRandomIntegers generates a list of randomly increasing integers with length "count".
// The list includes no duplicates.
func GenerateIncreasingRandomIntegers(count int) (list []int) {
	list = make([]int, count)

	for loop, index := 0, 0; index < len(list); loop++ {
		if r.Float32() > 0.5 {
			list[index] = loop
			index++
		}
	}
	return
}

// GenerateUniqueRandomIntegers generates a list of random unique integers.
func GenerateUniqueRandomIntegers(count int) (list sort.IntSlice) {
	list = make([]int, count)

	for i := range list {
		list[i] = i
	}

	// starting from the end, swap with a random smaller integer until done
	// (Fisher-Yates http://en.wikipedia.org/wiki/Fisher-Yates_shuffle)
	for i := count - 1; i > 0; i-- {
		// Intn is exclusive, Fisher-Yates says 0 <= j <= i
		rand := r.Intn(i + 1)
		list.Swap(i, rand)
	}

	return
}

// GenerateLimitedRandomIntegers generates a list of random integers that occur up 0..N times.
func GenerateLimitedRandomIntegers(count, occur int) (list sort.IntSlice) {
	// this is a horrifying way of making an always-growing list that's probably terrible
	list = make([]int, 0)

	for i := 0; i < count; i++ {
		for j := 0; j < r.Intn(occur+1); j++ {
			list = append(list, i)
		}
	}

	// starting from the end, swap with a random smaller integer until done
	// (Fisher-Yates http://en.wikipedia.org/wiki/Fisher-Yates_shuffle)
	for i := count - 1; i > 0; i-- {
		// Intn is exclusive, Fisher-Yates says 0 <= j <= i
		rand := r.Intn(i + 1)
		list.Swap(i, rand)
	}

	return
}

// GenerateRandomIntegers creates a list of 32bit integers with values up to bitcount bits
// 20 bitcount = 1<<20-1 max
func GenerateRandomIntegers(count int, bitcount uint32) (list []uint32, err error) {
	if bitcount > 32 {
		return nil, errors.New("limit must <= 32 bits")
	}
	list = make([]uint32, count)

	for i := 0; i < count; i++ {
		list[i] = r.Uint32() % uint32((1<<bitcount)-1)
	}
	return list, nil
}
