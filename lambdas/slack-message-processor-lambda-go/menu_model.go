package lambda

import (
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
)

// A type when evaluated on a group or option, is shown in the menu
type visible func() bool

func constVisible(show bool) visible {
	return func() bool {
		return show
	}
}

var (
	showOption = constVisible(true)
)

type adaptiveMenuOption struct {
	text  string
	value string
	show  visible
}

func menuOption(name, label string, show visible) adaptiveMenuOption {
	return adaptiveMenuOption{
		text:  label,
		value: name,
		show:  show,
	}
}

func staticMenuOption(name, label string) adaptiveMenuOption {
	return menuOption(name, label, showOption)
}

func resourceMenuOption(name, resourceKey string) adaptiveMenuOption {
	return staticMenuOption(name, RetrieveTemplate(resourceKey))
}

func (a adaptiveMenuOption) evaluate() bool {
	return a.show()
}

type adaptiveMenuOptionGroup struct {
	text    string
	options []adaptiveMenuOption
}

type adaptiveMenu struct {
	optionGroup adaptiveMenuOptionGroup
	show        visible
}

type allAdaptiveMenuOptions []adaptiveMenu

// Evaluate what options and groups to be shown in the menu
func (a allAdaptiveMenuOptions) evaluate() []ebm.MenuOptionGroup {
	var res []ebm.MenuOptionGroup
	for _, each := range a {
		if each.show() {
			var options []ebm.MenuOption
			group := each.optionGroup
			for _, ieach := range group.options {
				if ieach.evaluate() {
					// add only when visible is evaluated to true
					options = append(options, ebm.MenuOption{Text: ieach.text, Value: ieach.value})
				}
			}
			// show the group only where is at least one menu item in that group
			if len(options) > 0 {
				// add only when visible is evaluated to true
				res = append(res, ebm.MenuOptionGroup{Text: group.text, Options: options})
			}
		}
	}
	return res
}
