package ui

import (
	"fmt"
)

const (
	BulletBlackCircle   = "•" // • \u2022
	BulletWhiteCircle   = "◦" // ◦ \u25e6
	BulletBlackTriangle = "‣" // ‣ \u2023

	BulletTinyDot       = "․" // ․ \u2024

	BulletBlackSquare   = "▪" // ▪ black_small_square \u25aa :black_small_square:
	BulletWhiteSquare   = "▫" // ▫ white_small_square \u25ab :white_small_square:

	DefaultBullet       = BulletBlackCircle

)

// ListItem formats the text with bullet in front of it
func ListItem(text RichText) RichText {
	return DefaultBullet + " " + text + "\n"
}

// ListItems formats the given list of elements with bullets.
func ListItems(texts ...RichText) RichText {
	elements := RichTextMapToRichText(texts, ListItem)
	return RichTextConcat(elements...)
}

// NumberedList formats list with numbers
func NumberedList(numberSeparator string) func (texts ...RichText) RichText{
	return func (texts ...RichText) RichText {
		buffer := RichText("")
		for i, text := range texts {
			buffer += RichText(fmt.Sprintf("%v%s %s\n", i+1, numberSeparator, text))
		}
		return buffer
	}
}
