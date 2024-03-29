package dynamic_menu

import (
	// business_time "github.com/adaptiveteam/adaptive/business-time"
	// "github.com/adaptiveteam/adaptive/checks"
	"testing"

	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
)

func MenuForTesting() (menu DynamicMenuSpecification) {
	menu = NewAdaptiveDynamicMenu(
		NewGroupSpecification("Test Group #1",
			NewAdaptiveDynamicMenuSpecification(
				"test ID #1.1",
				"test text #1.1",
				"text description #1.1", true),
			NewAdaptiveDynamicMenuSpecification(
				"test ID #1.2",
				"test text #1.2",
				"text description #1.2", true),
		),
		NewGroupSpecification("Test Group #2",
			NewAdaptiveDynamicMenuSpecification(
				"test ID #2.1",
				"test text #2.1",
				"text description #2.1",
				true),
			NewAdaptiveDynamicMenuSpecification(
				"test ID #2.2",
				"test text #2.2",
				"text description #2.2",
				true),
		),
		NewGroupSpecification("Test Group #3",
			NewAdaptiveDynamicMenuSpecification(
				"test ID #3.1",
				"test text #3.1",
				"text description #3.1", true),
			NewAdaptiveDynamicMenuSpecification(
				"test ID #3.2",
				"test text #3.2",
				"text description #3.2", false),
		),
		NewGroupSpecification("Test Group #4",
			NewAdaptiveDynamicMenuSpecification(
				"test ID #4.1",
				"test text #4.1",
				"text description #4.1", false),
		),
	)
	return menu
}

func Test_ADM(t *testing.T) {

	desiredGroups := []string{
		"Test Group #1",
		"Test Group #2",
		"Test Group #3",
	}

	desiredOptions := map[string][]string{
		desiredGroups[0]: {
			"test ID #1.1",
			"test ID #1.2",
		},
		desiredGroups[1]: {
			"test ID #2.1",
			"test ID #2.2",
		},
		desiredGroups[2]: {
			"test ID #3.1",
		},
	}
	checkMenu(
		MenuForTesting(),
		desiredGroups,
		desiredOptions,
		t,
	)
}

func printStrings(v []string) (rv string) {
	for _, each := range v {
		rv = rv + "\t" + each + "\n"
	}
	return rv
}

func printOptions(v []ebm.MenuOption) (rv string) {
	for _, each := range v {
		rv = rv + "\t" + each.Value + "\n"
	}
	return rv
}

func checkMenu(
	menuSpec DynamicMenuSpecification,
	desiredGroups []string,
	desiredOptions map[string][]string,
	t *testing.T,
) {
	// "ctcreel", business_time.NewDate(2019, 1, 1)
	// result := menuConstructor()
	menu := menuSpec.Build()
	if len(menu) == len(desiredGroups) {
		for i, group := range desiredGroups {
			if len(desiredOptions[group]) == len(menu[i].Options) {
				if menu[i].Text == group {
					for j, option := range desiredOptions[group] {
						if menu[i].Options[j].Value != option {
							t.Errorf(
								"Expected option value %v but got %v",
								option,
								menu[i].Options[j].Value,
							)
						}
					}
				}
			} else {
				t.Errorf(
					"Expected length of options for %v to be %v but got %v.\nDesired:\n%v\nHave:\n%v",
					group,
					len(desiredOptions[group]),
					len(menu[i].Options),
					printStrings(desiredOptions[group]),
					printOptions(menu[i].Options),
				)
			}
		}
	} else {
		t.Errorf("Expected length of menu to be %v but got %v.",
			len(menu),
			len(desiredGroups),
		)
	}
}
