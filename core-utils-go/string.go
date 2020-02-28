package core_utils_go

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type AdaptiveUtils struct{}

var EmptyString = ""
var Underscore = "_"
var Asterisk = "*"

func ListContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Uuid() string {
	return uuid.New().String()
}

func TrimLower(ip string) string {
	return strings.TrimSpace(strings.ToLower(ip))
}

func UnorderedEqual(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[string]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	if len(diff) == 0 {
		return true
	}
	return false
}

func TextWrap(ip string, wrapper ...string) string {
	if ip != EmptyString {
		var res = ip
		for _, each := range wrapper {
			res = fmt.Sprintf("%s%s%s", each, res, each)
		}
		return res
	}
	return EmptyString
}

// Unique returns a unique subset of the string slice provided.
// Deprecated: use Distinct (it has more correct name)
func Unique(input []string) []string {
	return Distinct(input)
}

// ClipString clips string upto 'prefixLength' characters and adds suffix
func ClipString(str string, prefixLength int, suffix string) string {
	if len(str) < prefixLength {
		return str
	}
	return fmt.Sprintf("%s%s", str[0:prefixLength-len(suffix)], suffix)
}
