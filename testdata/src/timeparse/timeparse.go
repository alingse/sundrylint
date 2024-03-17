package timeparse

import (
	"fmt"
	"time"
)

const Format = "2006"

var FormatVar = "2006"

// nolint:unused
func callTimeParse() {
	var date = "2006"
	_, _ = time.Parse(date, time.DateOnly) // want `call func time.Parse may have incorrect parameters, potentially swapping the layout and value arguments.`
	_, _ = time.Parse(date, "2006")        // want `call func time.Parse may have incorrect parameters, potentially swapping the layout and value arguments.`
	_, _ = time.Parse(date, Format)        // want `call func time.Parse may have incorrect parameters, potentially swapping the layout and value arguments.`
	_, _ = time.Parse(date, FormatVar)
	timeFormatAny := any(time.DateOnly)
	_, _ = time.Parse(date, timeFormatAny.(string))
	var dates = []string{date}
	_, _ = time.Parse(date, dates[0])
	_, _ = time.Parse(date, fmt.Sprintf("%s-%s", "2006", "-1"))
	var strPtr = &date
	_, _ = time.Parse(date, *strPtr)

	var f = time.Parse
	_, _ = f(date, time.DateTime)
}
