package repeatargs

import (
	"time"
)

func GetTime(t int64) string {
	return time.Unix(t, 0).Format(time.DateTime)
}

func Call(startDate string, endDate string) string {
	return startDate + endDate
}

func Do(s int64, t int64) string {
	_ = s + t

	_ = Call(GetTime(s), GetTime(s)) // want `call the func with repeat args from a sub-func`

	var a A
	_ = Call(a.Name(), a.Name())

	return ""
}

type A struct {
}

func (A) Name() string {
	return "A"
}
