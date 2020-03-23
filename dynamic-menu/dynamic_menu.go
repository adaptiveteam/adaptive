package dynamic_menu

import (
	bt "github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/checks"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"log"
)

func NewAdaptiveDynamicMenu() (rv DynamicMenuSpecification) {
	rv = make(DynamicMenuSpecification, 0)
	return rv
}

func NewGroupSpecification(groupName string) (rv GroupSpecification) {
	rv = GroupSpecification{
		Group: ebm.MenuOptionGroup{
			Text: groupName,
		},
	}
	return rv
}

func NewAdaptiveDynamicMenuSpecification(
	optionID string,
	optionText string,
	optionDescription string,
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
			},
		}
	}
	return rv
}

func (option OptionSpecification) AddChecks(checks ...VisibilityCheck) OptionSpecification {
	option.Checks = append(option.Checks, checks...)
	return option
}

func (option OptionSpecification) AddOptionCheck(
	profile checks.CheckFunctionMap,
	name string,
	shouldBe bool,
) OptionSpecification {
	if profile[name] == nil {
		log.Panicln(name,"is not a valid check")
	}
	newOption := VisibilityCheck{
		Name:     name,
		Check:    profile[name],
		ShouldBe: shouldBe,
	}
	option.Checks = append(option.Checks, newOption)
	return option
}

// And adds a few option checks simultaneously
func (option OptionSpecification) And(
	profile checks.CheckFunctionMap,
	names ... string,
) OptionSpecification {
	for _, name := range names {
		option = option.AddOptionCheck(profile, name, true)
	}
	return option
}

func (gs GroupSpecification) AddGroupOption(option OptionSpecification) GroupSpecification {
	gs.Options = append(gs.Options, option)
	return gs
}

func (adm DynamicMenuSpecification) AddGroup(group GroupSpecification) (rv DynamicMenuSpecification) {
	rv = append(adm, group)
	return rv
}

// Build performs all checks and produces []ebm.MenuOptionGroup
func (adm DynamicMenuSpecification) Build(userID string, date bt.Date) (rv []ebm.MenuOptionGroup) {
	checkResults := adm.StripOutFunctions().Evaluate(userID, date)
	for _, currentGroup := range adm {
		groupOptions := make([]ebm.MenuOption, 0)
		for _, currentOption := range currentGroup.Options {
			optionIsVisible := true
			for i, check := range currentOption.Checks {
				optionIsVisible = optionIsVisible && (checkResults[check.Name] == currentOption.Checks[i].ShouldBe)
			}

			if optionIsVisible {
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
