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

	return Call(GetTime(s), GetTime(s)) // want `call the func with repeat args from a sub-func`
}
