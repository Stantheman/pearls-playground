package search

import (
	//"fmt"
	"github.com/Stantheman/pearls/helpers/binary"
	"github.com/Stantheman/pearls/helpers/random"
	"testing"
)

const (
	filename = "wow.txt"
	count    = 1048576
	bits     = 20
)

func TestMissing(t *testing.T) {
	// create a file with 1million 20-bit integers. since we can't make 20 bit integers
	// put them inside a 32-bit integer
	// the idea is at one point we created these, and now have
	// limited memory to work with
	ints, err := random.GenerateRandomIntegers(count, bits)
	if err != nil {
		t.Error(err)
	}
	if err := binary.MakeBinaryFile(filename, ints); err != nil {
		t.Error(err)
	}
	missing, err := Missing(filename, bits, count, 1, 0)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v is missing\n", missing)
	for i, _ := range ints {
		if ints[i] == missing {
			t.Fatal("it died here")
		}
	}
}
