package funcresultunused

import (
	"runtime/debug"
	"strconv"

	"funcresultunused/demo"
)

func FuncResultUnused() {
	var num int64 = 10
	strconv.FormatInt(num, 10) // want `func strconv.FormatInt return result is unused`

	_ = strconv.FormatBool(true)
	_ = map[string]string{
		"hello": strconv.FormatBool(true),
	}

	var s = new(XXStruct)
	s.FormatStructBool(true)
	// nolint:staticcheck
	strconv.FormatBool(true) // want `func strconv.FormatBool return result is unused`

	demo.FormatAny()   // want `func demo.FormatAny return result is unused`
	demo.FormatAny2(0) // want `func demo.FormatAny2 return result is unused`

	debug.SetGCPercent(100)
}

type XXStruct struct{}

func (x *XXStruct) FormatStructBool(t bool) string {
	return strconv.FormatBool(t)
}
