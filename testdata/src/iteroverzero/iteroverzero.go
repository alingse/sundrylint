package iteroverzero

func callAfterMake() int {
	var total int
	nums := make([]int, 0, 10)
	for _, i := range nums {
		total += i
	}
	return total
}

func callAfterMake2() int {
	var total int
	nums := make([]int, 0, 10)
	nums = append(nums, 1)
	for _, i := range nums {
		total += i
	}
	return total
}
