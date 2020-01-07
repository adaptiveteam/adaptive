package model

import (
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

type MenuOption struct {
	Text        string `json:"text"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

type MenuOptionGroup struct {
	Text    string       `json:"text"`
	Options []MenuOption `json:"options"`
}

// Option constructs menu option in the most widespread way.
func Option(name string, label ui.PlainText) MenuOption {
	return MenuOption{
		Text:  string(label),
		Value: name,
		Description: "",
	}
}

// SimpleOption constructs menu option for a simple case where value is equal to label
func SimpleOption(o ui.PlainText) MenuOption {
	return MenuOption{
		Text:  string(o),
		Value: string(o),
	}
}

// OptionGroup constructs menu option group from the given list of options.
func OptionGroup(title ui.PlainText, options ...MenuOption) MenuOptionGroup {
	return MenuOptionGroup{
		Text:    string(title),
		Options: options,
	}
}
