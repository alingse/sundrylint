package mustcompileout

import (
	"regexp"
)

const htmlTag = `\d+`

func ReplaceNumber(s string) string {
	r := regexp.MustCompile(htmlTag)
	return r.ReplaceAllString(s, "")
}
