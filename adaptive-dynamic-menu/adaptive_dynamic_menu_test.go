package adaptive_dynamic_menu

import (
	"fmt"
	acfn "github.com/adaptiveteam/adaptive/adaptive-check-function-names"
	bt "github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/checks"
	menu "github.com/adaptiveteam/adaptive/dynamic-menu"
	"github.com/adaptiveteam/adaptive/engagement-builder/model"
	"testing"
)

var AllTrueTestProfile = map[string]checks.CheckFunction{
	acfn.IDOsExistForMe:                            checks.ReturnsTrue,
	acfn.CanBeNudgedForIDO:                         checks.ReturnsTrue,
	acfn.StaleIDOsExistForMe:                       checks.ReturnsTrue,
	acfn.CompanyVisionExists:                       checks.ReturnsTrue,
	acfn.InStrategyCommunity:                       checks.ReturnsTrue,
	acfn.InitiativesExist:                          checks.ReturnsTrue,
	acfn.ObjectivesExistForMe:                      checks.ReturnsTrue,
	acfn.StaleObjectivesExistForMe:                 checks.ReturnsTrue,
	acfn.ObjectivesExistInMyCapabilityCommunities:  checks.ReturnsTrue,
	acfn.ObjectivesExist:                           checks.ReturnsTrue,
	acfn.InCapabilityCommunity:                     checks.ReturnsTrue,
	acfn.CapabilityCommunityExists:                 checks.ReturnsTrue,
	acfn.MultipleCapabilityCommunitiesExists:       checks.ReturnsTrue,
	acfn.InitiativesExistForMe:                     checks.ReturnsTrue,
	acfn.InitiativesExistInMyCapabilityCommunities: checks.ReturnsTrue,
	acfn.InitiativesExistInMyInitiativeCommunities: checks.ReturnsTrue,
	acfn.InitiativeCommunityExists:                 checks.ReturnsTrue,
	acfn.StaleInitiativesExistForMe:                checks.ReturnsTrue,
	acfn.InInitiativeCommunity:                     checks.ReturnsTrue,
	acfn.TeamValuesExist:                           checks.ReturnsTrue,
	acfn.InValuesCommunity:                         checks.ReturnsTrue,
	acfn.HolidaysExist:                             checks.ReturnsTrue,
	acfn.InHRCommunity:                             checks.ReturnsTrue,
	acfn.UndeliveredEngagementsExistForMe:          checks.ReturnsTrue,
	acfn.CollaborationReportExists:                 checks.ReturnsTrue,
	acfn.InLastMonthOfQuarter:                      checks.ReturnsTrue,
	acfn.UserSettingsExist:                         checks.ReturnsTrue,
	acfn.CoacheesExist:                             checks.ReturnsTrue,
	acfn.AdvocatesExist:                            checks.ReturnsTrue,
}

var AllFalseTestProfile = map[string]checks.CheckFunction{
	acfn.IDOsExistForMe:                            checks.ReturnsFalse,
	acfn.CanBeNudgedForIDO:                         checks.ReturnsFalse,
	acfn.StaleIDOsExistForMe:                       checks.ReturnsFalse,
	acfn.CompanyVisionExists:                       checks.ReturnsFalse,
	acfn.InStrategyCommunity:                       checks.ReturnsFalse,
	acfn.InitiativesExist:                          checks.ReturnsFalse,
	acfn.ObjectivesExistForMe:                      checks.ReturnsFalse,
	acfn.StaleObjectivesExistForMe:                 checks.ReturnsFalse,
	acfn.ObjectivesExistInMyCapabilityCommunities:  checks.ReturnsFalse,
	acfn.ObjectivesExist:                           checks.ReturnsFalse,
	acfn.InCapabilityCommunity:                     checks.ReturnsFalse,
	acfn.CapabilityCommunityExists:                 checks.ReturnsFalse,
	acfn.MultipleCapabilityCommunitiesExists:       checks.ReturnsFalse,
	acfn.InitiativesExistForMe:                     checks.ReturnsFalse,
	acfn.InitiativesExistInMyCapabilityCommunities: checks.ReturnsFalse,
	acfn.InitiativesExistInMyInitiativeCommunities: checks.ReturnsFalse,
	acfn.InitiativeCommunityExists:                 checks.ReturnsFalse,
	acfn.StaleInitiativesExistForMe:                checks.ReturnsFalse,
	acfn.InInitiativeCommunity:                     checks.ReturnsFalse,
	acfn.TeamValuesExist:                           checks.ReturnsFalse,
	acfn.InValuesCommunity:                         checks.ReturnsFalse,
	acfn.HolidaysExist:                             checks.ReturnsFalse,
	acfn.InHRCommunity:                             checks.ReturnsFalse,
	acfn.UndeliveredEngagementsExistForMe:          checks.ReturnsFalse,
	acfn.CollaborationReportExists:                 checks.ReturnsFalse,
	acfn.InLastMonthOfQuarter:                      checks.ReturnsFalse,
	acfn.UserSettingsExist:                         checks.ReturnsFalse,
	acfn.CoacheesExist:                             checks.ReturnsFalse,
	acfn.AdvocatesExist:                            checks.ReturnsFalse,
}

var IndividualContributor = map[string]checks.CheckFunction{
	acfn.IDOsExistForMe:                            checks.ReturnsFalse,
	acfn.CanBeNudgedForIDO:                         checks.ReturnsTrue,
	acfn.StaleIDOsExistForMe:                       checks.ReturnsFalse,
	acfn.CompanyVisionExists:                       checks.ReturnsTrue,
	acfn.InStrategyCommunity:                       checks.ReturnsFalse,
	acfn.InitiativesExist:                          checks.ReturnsFalse,
	acfn.ObjectivesExistForMe:                      checks.ReturnsFalse,
	acfn.StaleObjectivesExistForMe:                 checks.ReturnsFalse,
	acfn.ObjectivesExistInMyCapabilityCommunities:  checks.ReturnsTrue,
	acfn.ObjectivesExist:                           checks.ReturnsTrue,
	acfn.InCapabilityCommunity:                     checks.ReturnsFalse,
	acfn.CapabilityCommunityExists:                 checks.ReturnsFalse,
	acfn.MultipleCapabilityCommunitiesExists:       checks.ReturnsTrue,
	acfn.InitiativesExistForMe:                     checks.ReturnsFalse,
	acfn.InitiativesExistInMyCapabilityCommunities: checks.ReturnsFalse,
	acfn.InitiativesExistInMyInitiativeCommunities: checks.ReturnsTrue,
	acfn.InitiativeCommunityExists:                 checks.ReturnsTrue,
	acfn.StaleInitiativesExistForMe:                checks.ReturnsFalse,
	acfn.InInitiativeCommunity:                     checks.ReturnsTrue,
	acfn.TeamValuesExist:                           checks.ReturnsTrue,
	acfn.InValuesCommunity:                         checks.ReturnsFalse,
	acfn.HolidaysExist:                             checks.ReturnsTrue,
	acfn.InHRCommunity:                             checks.ReturnsFalse,
	acfn.UndeliveredEngagementsExistForMe:          checks.ReturnsFalse,
	acfn.CollaborationReportExists:                 checks.ReturnsFalse,
	acfn.InLastMonthOfQuarter:                      checks.ReturnsFalse,
	acfn.UserSettingsExist:                         checks.ReturnsTrue,
	acfn.CoacheesExist:                             checks.ReturnsTrue,
	acfn.AdvocatesExist:                            checks.ReturnsTrue,
}

var bindings = menu.FunctionBindings{
	"FetchEngagementsForMe":      "fetch_engagements_for_me",
	"StaleIDOsExistForMe":        "stale_idos_for_me",
	"StaleObjectivesExistForMe":  "stale_objectives_for_me",
	"StaleInitiativesExistForMe": "stale_initiatives_for_me",
	"AllIDOsForMe":               "all_idos_for_me",
	"AllIDOsForCoachees":         "all_idos_for_coachees",
	"AllObjectivesForMe":         "all_objectives_for_me",
	"AllInitiativesForMe":        "all_initiatives_for_me",
	"ProvideFeedback":            "provide_feedback",
	"RequestFeedback":            "request_feedback",
	"ViewEditObjectives":         "view_edit_objectives",
	"ViewCommunityObjectives":    "view_community_objectives",
	"ViewEditInitiatives":        "view_edit_initiatives",
	"ViewCommunityInitiatives":   "view_community_initiatives",
	"ViewValues":                 "view_values",
	"ViewEditValues":             "view_edit_values",
	"ViewCollaborationReport":    "view_collaboration_report",
	"ViewHolidays":               "view_holidays",
	"ViewEditHolidays":           "view_edit_holidays",
	"ViewScheduleNextQuarter":    "view_next_quarter_schedule",
	"ViewScheduleCurrentQuarter": "view_current_quarter_schedule",
	"ViewVision":                 "view_vision",
	"ViewEditVision":             "view_edit_vision",
	"CreateIDO":                  "create_ido",
	"CreateVision":               "create_vision",
	"CreateFinancialObjectives":  "create_financial_objective",
	"CreateCustomerObjectives":   "create_customer_objective",
	"CreateCapabilityObjectives": "create_capability_objective",
	"CreateCapabilityCommunity":  "createcapability_community",
	"CreateInitiatives":          "create_initiative",
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
			bindings["StaleIDOsExistForMe"],
			bindings["StaleObjectivesExistForMe"],
			bindings["StaleInitiativesExistForMe"],
			bindings["ProvideFeedback"],
			bindings["RequestFeedback"],
		},
		desiredGroups[1]: {
			bindings["AllIDOsForMe"],
			bindings["AllObjectivesForMe"],
			bindings["AllInitiativesForMe"],

			// No feedback here because we set the quarter check to true
		},
		desiredGroups[2]: {
			// User is in strategy community so can edit
			// User is in capability community so can edit
			// User is in initiative community so can edit
			// User is in HR community so can edit
			bindings["ViewEditVision"],
			bindings["ViewEditObjectives"],
			bindings["ViewEditInitiatives"],
			bindings["ViewEditValues"],
			bindings["ViewEditHolidays"],
			bindings["ViewScheduleCurrentQuarter"],
			bindings["ViewScheduleNextQuarter"],
			bindings["ViewCoachees"],
			bindings["ViewAdvocates"],
		},
		desiredGroups[3]: {
			// User is in strategy community so can create objectives
			// User is in capability community so create initiatives
			// User is in HR community so can create holidays
			bindings["CreateIDO"],
			bindings["CreateCapabilityObjectives"],
			bindings["CreateCapabilityCommunity"],
			bindings["CreateInitiatives"],
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
		AdaptiveDynamicMenu,
		AllTrueTestProfile,
		bindings,
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
	checkMenu(
		AdaptiveDynamicMenu,
		AllFalseTestProfile,
		bindings,
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
			bindings["CreateIDO"],
		},
		desiredGroups[1]: {
			bindings["ProvideFeedback"],
			bindings["RequestFeedback"],
		},
		desiredGroups[2]: {
			// User is in strategy community so can edit
			// User is in capability community so can edit
			// User is in initiative community so can edit
			// User is in HR community so can edit
			bindings["ViewVision"],
			bindings["ViewCommunityObjectives"],
			bindings["ViewCommunityInitiatives"],
			bindings["ViewValues"],
			bindings["ViewHolidays"],
			bindings["ViewScheduleCurrentQuarter"],
			bindings["ViewScheduleNextQuarter"],
			bindings["ViewCoachees"],
			bindings["ViewAdvocates"],
		},
		desiredGroups[3]: {
			bindings["CreateInitiatives"],
		},
		desiredGroups[4]: {
			bindings["UserSettings"],
		},
	}
	checkMenu(
		AdaptiveDynamicMenu,
		IndividualContributor,
		bindings,
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
	menuConstructor func(profile checks.CheckFunctionMap, bindings menu.FunctionBindings) menu.DynamicMenuSpecification,
	profile checks.CheckFunctionMap,
	bindings menu.FunctionBindings,
	desiredGroups []string,
	desiredOptions map[string][]string,
	t *testing.T,
) {
	result := menuConstructor(profile, bindings)
	newMenu := result.Build("ctcreel", bt.NewDate(2019, 1, 1))
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
