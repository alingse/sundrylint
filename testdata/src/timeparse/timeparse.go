package timeparse

import "time"

// nolint:unused
func callTimeParse() {
	var date = "2006"
	_, _ = time.Parse(date, time.DateOnly) // want `call func time.Parse may have incorrect parameters, potentially swapping the layout and value arguments.`
	_, _ = time.Parse(date, "2006")        // want `call func time.Parse may have incorrect parameters, potentially swapping the layout and value arguments.`
	var f = time.Parse

	_, _ = f(date, time.DateTime)
}
