package bin

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
