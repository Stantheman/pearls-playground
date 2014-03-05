/* Helpers is a chunk of functions that are useful to me in
testing code from programming pearls */
package random

import (
	"math/rand"
	"sort"
	"time"
)

// GenerateIncreasingRandomIntegers generates a list of randomly increasing integers with length count
// the list includes no duplicates
func GenerateIncreasingRandomIntegers(count int) (list []int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	list = make([]int, count)

	for loop, index := 0, 0; index < len(list); loop++ {
		if r.Float32() > 0.5 {
			list[index] = loop
			index++
		}
	}
	return
}

// GenerateUniqueRandomIntegers generates a list of random unique integers
func GenerateUniqueRandomIntegers(count int) (list sort.IntSlice) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
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
