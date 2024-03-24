package demo

import "strconv"

var N int64 = 1

func FormatAny() string {
	if N > 0 {
		return strconv.FormatInt(N, 10)
	}
	return "hello"
}

func FormatAny2(n int64) string {
	if n > 0 {
		return strconv.FormatInt(n, 10)
	}
	return "0"
}
