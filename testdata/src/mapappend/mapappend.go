package mapappend

func badAppend() {
	m := make(map[string][]int)
	m["b"] = append(m["a"], 1) // want "potential incorrect map append: using different keys for append and assignment"
}

func goodAppend() {
	m := make(map[string][]int)
	m["a"] = append(m["a"], 1) // 正确的用法，key相同
}
