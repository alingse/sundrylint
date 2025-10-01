package closureappend

func testClosureAppend() {
	// This should trigger the linter
	data := make([]int, 0, 10)

	for i := range 5 {
		go func() {
			data = append(data, i) // want `potential data race: append operation in closure captures loop variable or shared slice`
		}()
	}
}

func testClosureAppendWithFuture() {
	// This should trigger the linter
	var futures []func()
	data := make([]int, 0, 10)

	for i := range 5 {
		futures = append(futures, func() {
			data = append(data, i) // want `potential data race: append operation in closure captures loop variable or shared slice`
		})
	}
}

func testSafeClosureAppend() {
	// This should NOT trigger the linter
	for i := range 5 {
		go func(i int) {
			data := make([]int, 0, 1)
			data = append(data, i) // Safe: local slice and parameter
		}(i)
	}
}

func testSafeClosureNoLoop() {
	// This should NOT trigger the linter (no loop)
	data := make([]int, 0, 10)

	go func() {
		data = append(data, 1) // Safe: no loop variable capture
	}()
}

func testSafeClosureLocalSlice() {
	// This should NOT trigger the linter (local slice)
	for i := range 5 {
		go func() {
			data := make([]int, 0, 1)
			data = append(data, i) // Safe: local slice
		}()
	}
}