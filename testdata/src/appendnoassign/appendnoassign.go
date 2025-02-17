package appendnoassign

import (
	"io"
	"strconv"
)

func Do() {
	var b = make([]byte, 3)
	var c int64 = 1234567

	strconv.AppendInt(b, c, 10) // want `call strconv.AppendX but not keep func result`
	b = strconv.AppendBool(b, true)
	_ = b
}

func FormatInt(b []byte, c int64) []byte {
	return strconv.AppendInt(b, c, 10)
}

type X struct {
	D []byte
}

func WriteTo(b []byte, w io.Writer, c int64) {
	_, _ = w.Write(strconv.AppendInt(b, c, 10))

	_ = string(strconv.AppendBool(b, true))
	_ = X{
		D: strconv.AppendBool(b, false),
	}
}
