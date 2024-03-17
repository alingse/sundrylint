package iteroverzero

// nolint: unused
func callAfterMake() int {
	var total int
	nums := make([]int, 0)
	for _, i := range nums {
		total += i
	}
	return total
}
