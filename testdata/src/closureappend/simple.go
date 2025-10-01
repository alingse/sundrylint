package closureappend

import "sync"

// 简单测试：闭包中捕获循环变量并使用 append
func testSimpleLoopVariableCapture() {
	var wg sync.WaitGroup
	var results [][]int

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 错误：捕获循环变量 i 并使用 append
			results = append(results, []int{i}) // want "potential data race: append operation in closure captures loop variable or shared slice"
		}()
	}
	wg.Wait()
}