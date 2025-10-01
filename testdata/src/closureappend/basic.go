package closureappend

func testBasic() {
	var s []int

	// 在循环的闭包中使用 append
	for i := 0; i < 3; i++ {
		f := func() {
			s = append(s, 1) // want "potential data race: append operation in closure captures loop variable or shared slice"
		}
		f()
	}
}