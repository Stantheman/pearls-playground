package gcd

func EuclidGCD(a, b uint) uint {
	if b == 0 {
		return a
	}
	return EuclidGCD(b, a%b)
}

// from wikipedia, http://en.wikipedia.org/wiki/Binary_GCD_algorithm
func BinaryGCD(a, b uint) uint {
	// simple cases (termination)
	if a == b {
		return b
	}
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	// look for factors of 2
	if (^a & 1) != 0 {
		if (b & 1) != 0 { // b is odd
			return BinaryGCD(a>>1, b)
		} else {
			return (BinaryGCD(a>>1, b>>1) << 1)
		}
	}

	if (^b & 1) != 0 {
		return BinaryGCD(a, b>>1)
	}
	// reduce larger argument
	if a > b {
		return BinaryGCD((a-b)>>1, b)
	}
	return BinaryGCD((b-a)>>1, a)
}

// also adapted from wikipedia
func IterativeBinaryGCD(a, b uint) uint {
	/* GCD(0,b) == b; GCD(a,0) == a, GCD(0,0) == 0 */
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}

	/* Let shift := lg K, where K is the greatest power of 2
	   dividing both a and b. */
	var shift uint
	for shift = 0; ((a | b) & 1) == 0; shift++ {
		a >>= 1
		b >>= 1
	}

	// make a odd
	for (a & 1) == 0 {
		a >>= 1
	}
	/* From here on, a is always odd. */
	for {
		/* remove all factors of 2 in b -- they are not common */
		/*   note: b is not zero, so while will terminate */
		for (b & 1) == 0 { /* Loop X */
			b >>= 1
		}
		/* Now a and b are both odd. Swap if necessary so a <= b,
		   then set b = b - a (which is even). For bignums, the
		   swapping is just pointer movement, and the subtraction
		   can be done in-place. */
		if a > b {
			t := b
			b = a
			a = t
		} // Swap a and b.
		b = b - a // Here b >= a.
		if b == 0 {
			break
		}
	}
	/* restore common factors of 2 */
	return a << shift
}
