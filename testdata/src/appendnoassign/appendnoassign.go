package appendnoassign

import (
	"strconv"
)

func Do() {
	var b = make([]byte, 3)
	var c int64 = 1234567

	strconv.AppendInt(b, c, 10)
}
