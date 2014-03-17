// Trying Programming Pearls's first example
package main

import (
	"fmt"
)

func main() {
	question3()
}

func question3() {
	asterisks := "***************"

	fmt.Println("Question 3 asks about comparing the runtime of these functions versus the system sort")
	fmt.Println("This is the output from a debug run:")
	fmt.Println(asterisks + `
BenchmarkGoSort                   100                                          13252805  ns/op  555311  B/op  10033  allocs/op
BenchmarkNaiveSort                200                                          8809204   ns/op  171019  B/op  10011  allocs/op
BenchmarkLimitedSortOnePass       200                                          8634208   ns/op  171022  B/op  10011  allocs/op
BenchmarkBitSortOnePass           100                                          10217766  ns/op  103156  B/op  10327  allocs/op
BenchmarkBitSortPrimativeOnePass  200                                          9028865   ns/op  90356   B/op  10011  allocs/op
` + asterisks)
	fmt.Println("Each of these sorts, given enough memory, beat the native Go sort by taking advantage of the problem requirements\n")
	fmt.Println("The executable, compiled to run the NaiveSort is on par with the sort command as well:")
	fmt.Println(asterisks + `
		➜  pearls git:(master) ✗ time ../../../../bin/pearls
		../../../../bin/pearls  0.01s user 0.01s system 77% cpu 0.020 total
		➜  pearls git:(master) ✗ time sort -n /tmp/first_pearl_input.txt > /dev/null
		sort -n /tmp/first_pearl_input.txt > /dev/null  0.01s user 0.01s system 68% cpu 0.027 total
` + asterisks)
	fmt.Println("Introducing extra passes makes the speed decrease significantly and the memory usage increase.")
	fmt.Println("I'm pretty certain this is related to poor choices in IO:")
	fmt.Println(asterisks + `
BenchmarkGoSort                   100                                          13718540  ns/op  555328  B/op  10033   allocs/op
BenchmarkNaiveSort                200                                          9105642   ns/op  171016  B/op  10011   allocs/op
BenchmarkLimitedSort              50                                           33928727  ns/op  928900  B/op  100038  allocs/op
BenchmarkBitSort                  50                                           34806673  ns/op  861037  B/op  100358  allocs/op
BenchmarkBitSortPrimative         50                                           34283723  ns/op  848260  B/op  100038  allocs/op
` + asterisks)
}
