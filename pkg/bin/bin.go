package bin

import "math"

// NextPowerOf2 calculates the smallest power of 2 greater than or equal to n
func NextPowerOf2(n uint) uint {
	if n == 0 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}

// Pow2 math.Pow casted to uint
func Pow2(n uint) uint {
	return uint(math.Pow(2, float64(n)))
}
