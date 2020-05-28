package dynamic_menu

import (
	"github.com/adaptiveteam/adaptive/checks"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"log"
)

func NewAdaptiveDynamicMenu(groups ... GroupSpecification) (rv DynamicMenuSpecification) {
	rv = make(DynamicMenuSpecification, 0).AddGroups(groups...)
	return rv
}

func NewGroupSpecification(groupName string, options ... OptionSpecification) (rv GroupSpecification) {
	rv = GroupSpecification{
		Group: ebm.MenuOptionGroup{
			Text: groupName,
		},
		Options: options,
	}
	return rv
}

func NewAdaptiveDynamicMenuSpecification(
	optionID string,
	optionText string,
	optionDescription string,
	isActive bool,
) (rv OptionSpecification) {
	if len(optionText) == 0 {
		log.Panicln("optionText is empty")
	} else if len(optionID) == 0 {
		log.Panicln("optionID is empty for ",optionText)
	} else {
		rv = OptionSpecification{
			Option: ebm.MenuOption{
				Text:        optionText,
				Value:       optionID,
				Description: optionDescription,
				Active:      isActive,
			},
		}
	}
	return rv
}

func (gs GroupSpecification) AddGroupOptions(options ... OptionSpecification) GroupSpecification {
	gs.Options = append(gs.Options, options...)
	return gs
}

func (adm DynamicMenuSpecification) AddGroups(groups ... GroupSpecification) (rv DynamicMenuSpecification) {
	rv = append(adm, groups...)
	return rv
}

func (adm DynamicMenuSpecification) Build() (rv []ebm.MenuOptionGroup) {
	for _, currentGroup := range adm {
		groupOptions := make([]ebm.MenuOption, 0)
		for _, currentOption := range currentGroup.Options {
			if currentOption.Option.Active {
				groupOptions = append(
					groupOptions,
					ebm.MenuOption{
						Text:        currentOption.Option.Text,
						Value:       currentOption.Option.Value,
						Description: currentOption.Option.Description,
					},
				)
			}
		}
		if len(groupOptions) > 0 {
			rv = append(
				rv,
				ebm.MenuOptionGroup{
					Text:    currentGroup.Group.Text,
					Options: groupOptions,
				},
			)
		}
	}
	return rv
}

func (adm DynamicMenuSpecification) StripOutFunctions() (
	rv checks.CheckFunctionMap,
) {
	rv = make(checks.CheckFunctionMap, 0)
	// First, strip out all of the check functions and remove duplicates
	for _, currentGroup := range adm {
		for _, currentOption := range currentGroup.Options {
			for _, option := range currentOption.Checks {
				if option.Check == nil {
					log.Panicln(option,"not mapped to non-existent",option.Name,)
				} else {
					rv[option.Name] = option.Check
				}
			}
		}
	}
	return rv
}

type Profile struct {
	Map checks.CheckFunctionMap
}

func (p Profile)Check(name string) VisibilityCheck {
	if p.Map[name] == nil {
		log.Panicln(name,"is not a valid check")
	}
	return VisibilityCheck{
		Name:     name,
		Check:    p.Map[name],
		ShouldBe: true,
	}
}

func (p Profile)Checks(names ... string) VisibilityCheck {
	checkName := ""
	check := checks.True
	for _, name := range names {
		if checkName != "" { checkName = checkName + " && " }
		checkName = checkName + name
		check = check.And(p.Map.GetUnsafe(name))
	}
	return VisibilityCheck{
		Name:     checkName,
		Check:    check,
		ShouldBe: true,
	}
}

func (c VisibilityCheck)Inv() VisibilityCheck {
	c.ShouldBe = !c.ShouldBe
	return c
}
