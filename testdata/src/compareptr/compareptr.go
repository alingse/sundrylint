package compareptr

import "fmt"

func Demo() {
	var a int64 = 2
	var b int64 = 1
	var c = &a
	var d = &b
	if c == d { // want "comparing pointers with == or != can be error-prone"
		fmt.Println("ok")
	}
}
