package adaptive_dynamic_menu

import (
	"fmt"
	"testing"

	"github.com/adaptiveteam/adaptive/adaptive-checks"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"

	// acfn "github.com/adaptiveteam/adaptive/adaptive-check-function-names"
	// bt "github.com/adaptiveteam/adaptive/business-time"
	// "github.com/adaptiveteam/adaptive/checks"
	menu "github.com/adaptiveteam/adaptive/dynamic-menu"
	"github.com/adaptiveteam/adaptive/engagement-builder/model"
)

var bindings = menu.FunctionBindings{
	"FetchEngagementsForMe":      "fetch_engagements_for_me",
	"StaleObjectivesExistForMe":  "stale_objectives_for_me",
	"StaleInitiativesExistForMe": "stale_initiatives_for_me",
	"AllIDOsForCoachees":         "all_idos_for_coachees",
	"AllObjectivesForMe":         "all_objectives_for_me",
	"AllInitiativesForMe":        "all_initiatives_for_me",
	"ProvideFeedback":            "provide_feedback",
	"RequestFeedback":            "request_feedback",
	"ViewValues":                 "view_values",
	"ViewEditValues":             "view_edit_values",
	"ViewCollaborationReport":    "view_collaboration_report",
	"ViewHolidays":               "view_holidays",
	"ViewEditHolidays":           "view_edit_holidays",
	"ViewScheduleNextQuarter":    "view_next_quarter_schedule",
	"ViewScheduleCurrentQuarter": "view_current_quarter_schedule",
	"ViewVision":                 "view_vision",
	"ViewEditVision":             "view_edit_vision",
	"CreateVision":               "create_vision",
	"CreateFinancialObjectives":  "create_financial_objective",
	"CreateCustomerObjectives":   "create_customer_objective",
	"CreateCapabilityCommunity":  "createcapability_community",
	"CreateInitiativeCommunity":  "create_initiative_community",
	"CreateValues":               "create_value",
	"CreateHolidays":             "create_holidays",
	"AssignCapabilityObjective":  "assign_capability_objective",
	"UserSettings":               "update_settings",
	"ViewCoachees":               "view_coachees",
	"ViewAdvocates":              "view_advocates",
	"StrategyPerformanceReport":  "StrategyPerformanceReport",
	"IDOPerformanceReport":       "IDOPerformanceReport",
}

func Test_AllTrue(t *testing.T) {
	desiredGroups := []string{
		"Urgent Responsibilities",
		"Responsibilities",
		"View",
		"Create",
		"Assign",
		"Settings",
		"Reports",
	}

	desiredOptions := map[string][]string{
		desiredGroups[0]: {
			bindings["FetchEngagementsForMe"],
			user.StaleIDOsForMe,
			bindings["StaleObjectivesExistForMe"],
			bindings["StaleInitiativesExistForMe"],
			bindings["ProvideFeedback"],
			bindings["RequestFeedback"],
		},
		desiredGroups[1]: {
			user.ViewObjectives,
			bindings["AllObjectivesForMe"],
			bindings["AllInitiativesForMe"],

			// No feedback here because we set the quarter check to true
		},
		desiredGroups[2]: {
			// User is in strategy community so can edit
			// User is in objective community so can edit
			// User is in initiative community so can edit
			// User is in HR community so can edit
			bindings["ViewEditVision"],
			strategy.ViewStrategyObjectives,
			strategy.ViewCapabilityCommunityInitiatives,
			bindings["ViewEditValues"],
			bindings["ViewEditHolidays"],
			bindings["ViewScheduleCurrentQuarter"],
			bindings["ViewScheduleNextQuarter"],
			bindings["ViewCoachees"],
			bindings["ViewAdvocates"],
		},
		desiredGroups[3]: {
			// User is in strategy community so can create objectives
			// User is in objective community so create initiatives
			// User is in HR community so can create holidays
			objectives.CreateIDONow,
			strategy.CreateStrategyObjective,
			bindings["CreateCapabilityCommunity"],
			strategy.CreateInitiative,
			bindings["CreateInitiativeCommunity"],
			bindings["CreateValues"],
			bindings["CreateHolidays"],
		},
		desiredGroups[4]: {
			bindings["AssignCapabilityObjective"],
		},
		desiredGroups[5]: {
			bindings["UserSettings"],
		},
		desiredGroups[6]: {
			bindings["ViewCollaborationReport"],
			bindings["StrategyPerformanceReport"],
			bindings["IDOPerformanceReport"],
		},
	}
	checkMenu(
		AdaptiveDynamicMenu(adaptive_checks.SomeTrueAndSomeFalseTestProfile, bindings),
		desiredGroups,
		desiredOptions,
		t,
	)
}

func Test_AllFalse(t *testing.T) {
	desiredGroups := []string{
		"View",
		"Settings",
	}

	desiredOptions := map[string][]string{
		desiredGroups[0]: {
			bindings["ViewScheduleCurrentQuarter"],
			bindings["ViewScheduleNextQuarter"],
			// bindings["ViewCoachees"],
		},
		desiredGroups[1]: {
			bindings["UserSettings"],
		},
	}
	// "ctcreel", bt.NewDate(2019, 1, 1)
	checkMenu(
		AdaptiveDynamicMenu(adaptive_checks.AllFalseTestProfile, bindings),
		desiredGroups,
		desiredOptions,
		t,
	)
}

func Test_IndividualContributor(t *testing.T) {
	desiredGroups := []string{
		"Urgent Responsibilities",
		"Responsibilities",
		"View",
		"Create",
		"Settings",
	}

	desiredOptions := map[string][]string{
		desiredGroups[0]: {
			objectives.CreateIDONow,
		},
		desiredGroups[1]: {
			bindings["ProvideFeedback"],
			bindings["RequestFeedback"],
		},
		desiredGroups[2]: {
			// User is in strategy community so can edit
			// User is in objective community so can edit
			// User is in initiative community so can edit
			// User is in HR community so can edit
			bindings["ViewVision"],
			strategy.ViewCapabilityCommunityObjectives,
			strategy.ViewInitiativeCommunityInitiatives,
			bindings["ViewValues"],
			bindings["ViewHolidays"],
			bindings["ViewScheduleCurrentQuarter"],
			bindings["ViewScheduleNextQuarter"],
			bindings["ViewCoachees"],
			bindings["ViewAdvocates"],
		},
		desiredGroups[3]: {
			strategy.CreateInitiative,
		},
		desiredGroups[4]: {
			bindings["UserSettings"],
		},
	}
	checkMenu(
		AdaptiveDynamicMenu(adaptive_checks.IndividualContributor, bindings),
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

func printOptions(v []model.MenuOption) (rv string) {
	for _, each := range v {
		rv = rv + "\t" + each.Value + "\n"
	}
	return rv
}

func checkMenu(
	result menu.DynamicMenuSpecification,
	desiredGroups []string,
	desiredOptions map[string][]string,
	t *testing.T,
) {
	newMenu := result.Build()
	if len(newMenu) == len(desiredGroups) {
		for i, group := range desiredGroups {
			if len(desiredOptions[group]) == len(newMenu[i].Options) {
				if newMenu[i].Text == group {
					for j, option := range desiredOptions[group] {
						if newMenu[i].Options[j].Value != option {
							t.Errorf(
								"Expected option value %v but got %v",
								option,
								newMenu[i].Options[j].Value,
							)
						}
					}
				}
			} else {
				t.Errorf(
					"Expected length of options for %v to be %v but got %v.\nDesired:\n%v\nHave:\n%v",
					group,
					len(desiredOptions[group]),
					len(newMenu[i].Options),
					printStrings(desiredOptions[group]),
					printOptions(newMenu[i].Options),
				)
			}
		}
	} else {
		fmt.Println(newMenu)
		fmt.Println(desiredGroups)
		t.Errorf("Expected length of menu to be %v but got %v.",
			len(desiredGroups),
			len(newMenu),
		)
	}
}
