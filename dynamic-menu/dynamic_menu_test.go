package dynamic_menu

import (
	business_time "github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/checks"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"testing"
)

func MenuForTesting(profile checks.CheckFunctionMap) (menu DynamicMenuSpecification) {
	menu = NewAdaptiveDynamicMenu()
	menu = menu.AddGroup(
		NewGroupSpecification("Test Group #1").
			AddGroupOption(
				NewAdaptiveDynamicMenuSpecification(
					"test ID #1.1",
					"test text #1.1",
					"text description #1.1").
					AddOptionCheck(profile, "ReturnsTrue", true).
					AddOptionCheck(profile, "ReturnsTrue", true),
			).AddGroupOption(
			NewAdaptiveDynamicMenuSpecification(
				"test ID #1.2",
				"test text #1.2",
				"text description #1.2").
				AddOptionCheck(profile, "ReturnsTrue", true).
				AddOptionCheck(profile, "ReturnsTrue", true),
		),
	).AddGroup(
		NewGroupSpecification("Test Group #2").
			AddGroupOption(
				NewAdaptiveDynamicMenuSpecification(
					"test ID #2.1",
					"test text #2.1",
					"text description #2.1").
					AddOptionCheck(profile, "ReturnsTrue", true).
					AddOptionCheck(profile, "ReturnsTrue", true),
			).AddGroupOption(
			NewAdaptiveDynamicMenuSpecification(
				"test ID #2.2",
				"test text #2.2",
				"text description #2.2").
				AddOptionCheck(profile, "ReturnsTrue", true).
				AddOptionCheck(profile, "ReturnsTrue", true),
		),
	).AddGroup(
		NewGroupSpecification("Test Group #3").
			AddGroupOption(
				NewAdaptiveDynamicMenuSpecification(
					"test ID #3.1",
					"test text #3.1",
					"text description #3.1").
					AddOptionCheck(profile, "ReturnsTrue", true).
					AddOptionCheck(profile, "ReturnsTrue", true),
			).AddGroupOption(
			NewAdaptiveDynamicMenuSpecification(
				"test ID #3.2",
				"test text #3.2",
				"text description #3.2").
				AddOptionCheck(profile, "ReturnsTrue", true).
				AddOptionCheck(profile, "ReturnsTrue", false),
		),
	).AddGroup(
		NewGroupSpecification("Test Group #4").
			AddGroupOption(
				NewAdaptiveDynamicMenuSpecification(
					"test ID #4.1",
					"test text #4.1",
					"text description #4.1").
					AddOptionCheck(profile, "ReturnsTrue", false).
					AddOptionCheck(profile, "ReturnsTrue", false),
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
		MenuForTesting,
		checks.SimpleTestProfile,
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
	menuConstructor func(profile checks.CheckFunctionMap) DynamicMenuSpecification,
	profile checks.CheckFunctionMap,
	desiredGroups []string,
	desiredOptions map[string][]string,
	t *testing.T,
) {
	result := menuConstructor(profile)
	menu := result.Build("ctcreel", business_time.NewDate(2019, 1, 1))
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
