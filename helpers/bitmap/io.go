package bitmap

import (
	"fmt"
	"os"
)

/* sortSetup is a helper sub to do the grunt work before the sort begins.
Simple error checking, returns two filehandles or possibly an error. Caller
must check error and is responsible for calling close or deferring */
func sortSetup(input_fn, output_fn string, length, avail int) (in, out *os.File, err error) {
	if length < 0 {
		return nil, nil, fmt.Errorf("Length must be greater than 0: %v", length)
	}

	if avail < 0 {
		return nil, nil, fmt.Errorf("Avail must be greater than 0: %v\n", avail)
	}

	// open the input file for reading
	in, err = os.Open(input_fn)
	if err != nil {
		return nil, nil, err
	}

	// same for output
	out, err = os.Create(output_fn)
	if err != nil {
		return nil, nil, err
	}

	return
}
