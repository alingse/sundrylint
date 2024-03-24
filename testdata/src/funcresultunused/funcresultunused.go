package funcresultunused

import (
	"strconv"

	"funcresultunused/demo"
)

func FuncResultUnused() {
	var num int64 = 10
	strconv.FormatInt(num, 10) // want `func result unused`

	_ = strconv.FormatBool(true)
	_ = map[string]string{
		"hello": strconv.FormatBool(true),
	}

	var s = new(XXStruct)
	s.FormatStructBool(true)
	// nolint:staticcheck
	strconv.FormatBool(true) // want `func result unused`

	demo.FormatAny()   // want `func result unused`
	demo.FormatAny2(0) // want `func result unused`
}

type XXStruct struct{}

func (x *XXStruct) FormatStructBool(t bool) string {
	return strconv.FormatBool(t)
}
