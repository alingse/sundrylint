package closureappend

// 最简单的测试用例
func testDebug() {
	var results []int

	// 在循环的闭包中使用 append
	for i := 0; i < 3; i++ {
		f := func() {
			results = append(results, 1) // want "potential data race: append operation in closure captures loop variable or shared slice"
		}
		f()
	}
}