package appendnoassign

import (
	"strconv"
)

func Do() {
	var b = make([]byte, 3)
	var c int64 = 1234567

	strconv.AppendInt(b, c, 10) // want `call strconv.AppendX but not keep func result`
	b = strconv.AppendBool(b, true)
	_ = b
}
