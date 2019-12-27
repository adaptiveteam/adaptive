package ui

import (
	"fmt"
	"strings"
)

// Bold wraps the text in * (asterisk)
func (s RichText)Bold() RichText {
	return s.Wrap("*")
}

// Italics wraps the text in _ (underscore)
func (s RichText)Italics() RichText {
	return s.Wrap("_")
}

// Monospaced wraps the text in '`' (backticks)
func (s RichText)Monospaced() RichText {
	return s.Wrap("`")
}

// Wrap wraps the text in wrapper
func (s RichText)Wrap(wrapper string) RichText {
	if(s == RichText("")) {
		return ""
	}
	return RichText(wrapper) + s + RichText(wrapper)
}

// TextWrap wraps text in a few wrappers.
func (s RichText)TextWrap(wrappers ...string) RichText {
	if(s == RichText("")) {
		return ""
	}
	var res = s
	for _, wrapper := range wrappers {
		res = RichText(fmt.Sprintf("%s%s%s", wrapper, res, wrapper))
	}
	return res
}

// UserTag formats user id in a way that will show Display Name of the user.
func UserTag(userID string) RichText {
	return RichText("<@" + userID + ">")
}

// ToRichText converts plain text to RichText
func (t PlainText)ToRichText() RichText {
	return RichText(t)
}

// Sprintf formats text using provided template
// internally it uses ordinary fmt.Sprintf
func Sprintf(format RichText, a ...interface{}) RichText {
	return RichText(fmt.Sprintf(string(format), a...))
}

// Elipsis limits text length to the provided maximum
func (t PlainText)Elipsis(maxLength int) PlainText {
	if len(t) < maxLength {
		return t
	}
	return t[:maxLength-3] + "..."
}
// Join concatenates elements of `a` placing `sep` in between.
func Join(a []RichText, sep RichText) RichText {
	s := make([]string, len(a))
	for i := 0; i < len(a); i++ {
		s[i] = string(a[i])
	}
	return RichText(strings.Join(s, string(sep)))
}
