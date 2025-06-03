package mapappend

func badAppend() {
	m := make(map[string][]int)
	m["b"] = append(m["a"], 1) // want "potential incorrect map append: using different keys for append and assignment"
}

func goodAppend() {
	m := make(map[string][]int)
	m["a"] = append(m["a"], 1) // 正确的用法，key相同
}

type TimeType int

const (
	Month TimeType = 1
	Day   TimeType = 1
)

func badAppend2() {
	m := make(map[TimeType][]int)
	m[Month] = append(m[Day], 1) // want `potential incorrect map append: using different keys for append and assignment`
}
