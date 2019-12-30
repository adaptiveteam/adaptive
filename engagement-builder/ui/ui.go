package ui

import "bytes"

// PlainText is a string that will be shown to user. Does not support formatting
type PlainText string
// RichText is a string that will be shown to user. Supports formatting
type RichText string

// RichTextMapToRichText applies function to each element of the input slice
func RichTextMapToRichText(vs []RichText, f func(RichText) RichText) []RichText {
    vsm := make([]RichText, len(vs))
    for i, v := range vs {
        vsm[i] = f(v)
    }
    return vsm
}

// RichTextConcat concatenates all strings
func RichTextConcat(lst ...RichText) RichText {
	var buffer bytes.Buffer

    for i := 0; i < len(lst); i++ {
        buffer.WriteString(string(lst[i]))
    }

    return RichText(buffer.String())
}
