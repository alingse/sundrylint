package mustcompileout

import (
	"regexp"
)

const htmlTag = `\d+`

var _ = regexp.MustCompile(htmlTag)

func ReplaceNumber(s string) string {
	r := regexp.MustCompile(htmlTag) // want `call regexp.MustCompile with constant should be moved out of func`
	return r.ReplaceAllString(s, "")
}

func Eval(rg, s string) string {
	r := regexp.MustCompile(rg)
	return r.ReplaceAllString(s, "")
}
