package rangeappendall

// nolint: unused
func collectBigger(ns []int, k int) []int {
	rs := make([]int, 0)
	for _, n := range ns {
		if n > k {
			rs = append(rs, ns...) // want `append all its data while range it`
		}
	}
	return rs
}

// nolint: unused
func collectLong(ns []string, k int) []string {
	rs := make([]string, 0)
	for _, n := range ns {
		if len(n) > k {
			rs = append(rs, ns...) // want `append all its data while range it`
		}
	}
	return rs
}

// nolint: unused
func collectLongCorrect(ns []string, k int) []string {
	rs := make([]string, 0)
	n := ns[0]
	rs = append(rs, n)
	for _, n := range ns {
		if len(n) > k {
			rs = append(rs, n)
		}
	}
	rs = append(rs, ns...)
	return rs
}

func CollectWithNil(ns []string, k int) bool {
	rs := make([]string, 0)
	rs2 := make([]string, 0)
	rs3 := make([]string, 0)
	for _, n := range ns {
		if len(n) > k {
			rs = append([]string{n}, ns...)
			rs2 = append([]string{}, ns...)
			rs3 = append([]string(nil), ns...)
		}
	}
	if len(rs) > 0 || len(rs2) > 0 || len(rs3) > 0 {
		return true
	}
	return false
}
