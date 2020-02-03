package adaptive_dynamic_menu

import (
	acfn "github.com/adaptiveteam/adaptive/adaptive-check-function-names"
	"github.com/adaptiveteam/adaptive/checks"
	menu "github.com/adaptiveteam/adaptive/dynamic-menu"
)

func AdaptiveDynamicMenu(profile checks.CheckFunctionMap, bindings menu.FunctionBindings) (adm menu.DynamicMenuSpecification) {
	p := menu.Profile{Map: profile}
	adm = menu.NewAdaptiveDynamicMenu()
	adm = adm.AddGroup(
		menu.NewGroupSpecification("Urgent Responsibilities").
			AddGroupOption(
				// Enables the user to create the company vision
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["CreateVision"],
					"Create Vision",
					"").
					AddOptionCheck(profile, acfn.CompanyVisionExists, false).
					AddOptionCheck(profile, acfn.InStrategyCommunity, true),
			).AddGroupOption(
			// This fetches any undelivered engagements
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["FetchEngagementsForMe"],
				"What do I have right now?",
				"").
				AddOptionCheck(profile, acfn.UndeliveredEngagementsOrPostponedEventsExistForMe, true),
		).AddGroupOption(
			// This fetches any IDO's not updated in the last 7 days
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["StaleIDOsExistForMe"],
				"Update IDO's",
				"").
				AddOptionCheck(profile, acfn.StaleIDOsExistForMe, true),
		).AddGroupOption(
			// This fetches any Objectives not updated in the last 7 days
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["StaleObjectivesExistForMe"],
				"Update Objectives",
				"").
				AddOptionCheck(profile, acfn.StaleObjectivesExistForMe, true),
		).AddGroupOption(
			// This fetches any Initiatives not updated in the last 7 days
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["StaleInitiativesExistForMe"],
				"Update Initiatives",
				"").
				AddOptionCheck(profile, acfn.StaleInitiativesExistForMe, true),
		).AddGroupOption(
			// Enables the user to create an IDO
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateIDO"],
				"Create IDO",
				"").
				AddOptionCheck(profile, acfn.CanBeNudgedForIDO, true).
				AddOptionCheck(profile, acfn.IDOsExistForMe, false).
				AddOptionCheck(profile, acfn.TeamValuesExist, true),
		).AddGroupOption(
			// This enables the user to post feedback to another user
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ProvideFeedback"],
				"Provide Feedback",
				"").
				AddOptionCheck(profile, acfn.InLastMonthOfQuarter, true).
				AddOptionCheck(profile, acfn.TeamValuesExist, true),
		).AddGroupOption(
			// This enables the user to request feedback from another user
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["RequestFeedback"],
				"Request Feedback",
				"").
				AddOptionCheck(profile, acfn.InLastMonthOfQuarter, true).
				AddOptionCheck(profile, acfn.TeamValuesExist, true),
		),
	).AddGroup(
		menu.NewGroupSpecification("Responsibilities").
			AddGroupOption(
				// This fetches all IDO's
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["AllIDOsForMe"],
					"All IDO's",
					"").
					AddOptionCheck(profile, acfn.IDOsExistForMe, true),
			).AddGroupOption(
			// This fetches all Objectives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["AllObjectivesForMe"],
				"All Objectives",
				"").
				AddOptionCheck(profile, acfn.ObjectivesExistForMe, true),
		).AddGroupOption(
			// This fetches all Initiatives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["AllInitiativesForMe"],
				"All Initiatives",
				"").
				AddOptionCheck(profile, acfn.InitiativesExistForMe, true),
		).AddGroupOption(
			// This enables the user to post feedback to another user
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ProvideFeedback"],
				"Provide Feedback",
				"").
				AddOptionCheck(profile, acfn.InLastMonthOfQuarter, false).
				AddOptionCheck(profile, acfn.TeamValuesExist, true),
		).AddGroupOption(
			// This enables the user to request feedback from another user
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["RequestFeedback"],
				"Request Feedback",
				"").
				AddOptionCheck(profile, acfn.InLastMonthOfQuarter, false).
				AddOptionCheck(profile, acfn.TeamValuesExist, true),
		),
	).AddGroup(
		menu.NewGroupSpecification("View").
			AddGroupOption(
				// Presents the company vision
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["ViewVision"],
					"Vision",
					"").
					AddOptionCheck(profile, acfn.CompanyVisionExists, true).
					AddOptionCheck(profile, acfn.InStrategyCommunity, false),
			).AddGroupOption(
			// Presents the company vision
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewEditVision"],
				"Vision",
				"").
				AddOptionCheck(profile, acfn.CompanyVisionExists, true).
				AddOptionCheck(profile, acfn.InStrategyCommunity, true),
		).AddGroupOption(
			// Enables the user to see all of the objectives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewEditObjectives"],
				"Objectives",
				"").
				AddOptionCheck(profile, acfn.ObjectivesExist, true).
				AddOptionCheck(profile, acfn.InStrategyCommunity, true),
		).AddGroupOption(
			// Enables the user to see all of the objectives in their Capability & Initiative Communities
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewCommunityObjectives"],
				"Objectives",
				"").
				AddOptionCheck(profile, acfn.ObjectivesExistInMyCapabilityCommunities, true).
				AddOptionCheck(profile, acfn.InStrategyCommunity, false),
		).AddGroupOption(
			// Enables the user to see all of the initiatives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewCommunityInitiatives"],
				"Initiatives",
				"").
				AddOptionCheck(profile, acfn.InStrategyCommunity, false).
				AddOptionCheck(profile, acfn.InitiativesExistInMyInitiativeCommunities, true).
				AddOptionCheck(profile, acfn.InCapabilityCommunity, false),
		).AddGroupOption(
			// Enables the user to see and edit all of the Initiatives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewEditInitiatives"],
				"Initiatives",
				"").
				AddOptionCheck(profile, acfn.InStrategyCommunity, false).
				AddOptionCheck(profile, acfn.InitiativesExistInMyCapabilityCommunities, true).
				AddOptionCheck(profile, acfn.InCapabilityCommunity, true),
		).AddGroupOption(
			// Enables the user to see and edit all of the Initiatives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewEditInitiatives"],
				"Initiatives",
				"").
				AddOptionCheck(profile, acfn.InStrategyCommunity, true).
				AddOptionCheck(profile, acfn.InitiativesExist, true),
		).AddGroupOption(
			// Enables the user to see all of the team values
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewEditValues"],
				"Competencies",
				"").
				AddOptionCheck(profile, acfn.TeamValuesExist, true).
				AddOptionCheck(profile, acfn.InValuesCommunity, true),
		).AddGroupOption(
			// Enables the user to see all of the team values
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewValues"],
				"Competencies",
				"").
				AddOptionCheck(profile, acfn.TeamValuesExist, true).
				AddOptionCheck(profile, acfn.InValuesCommunity, false),
		).AddGroupOption(
			// Enables the user to see all of the company holidays
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewHolidays"],
				"Holidays",
				"").
				AddOptionCheck(profile, acfn.HolidaysExist, true).
				AddOptionCheck(profile, acfn.InHRCommunity, false),
		).AddGroupOption(
			// Enables the user to see all of the company holidays
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewEditHolidays"],
				"Holidays",
				"").
				AddOptionCheck(profile, acfn.HolidaysExist, true).
				AddOptionCheck(profile, acfn.InHRCommunity, true),
		).AddGroupOption(
			// Enables the user to see all of the scheduled events for the current quarter
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewScheduleCurrentQuarter"],
				"Current Quarter Events",
				""),
		).AddGroupOption(
			// Enables the user to see all of the scheduled events for the next quarter
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewScheduleNextQuarter"],
				"Next Quarter Events",
				""),
		).AddGroupOption(
			// Enables the user to see all of the IDOs they are coaching
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewCoachees"],
				"Coachee IDOs",
				"").
				AddOptionCheck(profile, acfn.CoacheesExist, true),
		).AddGroupOption(
			// Enables the user to see all of the advocates for strategic elements aligned to you
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewAdvocates"],
				"Advocates",
				"").
				AddOptionCheck(profile, acfn.AdvocatesExist, true),
		),
	).AddGroup(
		menu.NewGroupSpecification("Create").
			AddGroupOption(
				// Enables a user to create an IDO
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["CreateIDO"],
					"IDO",
					"").
					AddOptionCheck(profile, acfn.CanBeNudgedForIDO, true).
					AddOptionCheck(profile, acfn.TeamValuesExist, true).
					AddOptionCheck(profile, acfn.IDOsExistForMe, true),
			).AddGroupOption(
			// Enables a non user/initiative community user to create an IDO
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateIDO"],
				"IDO",
				"").
				AddOptionCheck(profile, acfn.CanBeNudgedForIDO, false).
				AddOptionCheck(profile, acfn.TeamValuesExist, true),
		).AddGroupOption(
			// Enables the user to create objectives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateCapabilityObjectives"],
				"Objectives",
				"").
				AddOptionCheck(profile, acfn.CompanyVisionExists, true).
				AddOptionCheck(profile, acfn.InStrategyCommunity, true).
				AddOptionCheck(profile, acfn.CapabilityCommunityExists, true),
		).AddGroupOption(
			// Enables the user to create Objective Communities
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateCapabilityCommunity"],
				"Objective Communities",
				"").
				AddOptionCheck(profile, acfn.CompanyVisionExists, true).
				AddOptionCheck(profile, acfn.InStrategyCommunity, true),
//***************************************************************************************************//
// The next four entries enable people who are either in a capability community to
// create an initiative community or initiative if they are in the strategy community
// or in an associated capability community.  The first set of two group options
// enables people who are in a capability community to create an initiative community
// (if there is a capability community) or an initiative
// (if there is an initiative community). The group does the same but for people who
// are in the strategy community.  To prevent duplicates we in the second group we
//  explicitely exclude the conditions from the first group.
// 
// TODO: Actually this should be implemented using boolean logic (and, or, not ...)
/////////////////////////////////////////////////////////////////////////////////////////////////////////
		).AddGroupOption(
			// Enables the user to create Initiatives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateInitiatives"],
				"Initiatives",
				"").
				AddOptionCheck(profile, acfn.InStrategyCommunity, false).
				AddOptionCheck(profile, acfn.ObjectivesExistInMyCapabilityCommunities, true).
				AddOptionCheck(profile, acfn.InInitiativeCommunity, true),
			).AddGroupOption(
				// Enables users from the strategy community to create Initiatives
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["CreateInitiatives"],
					"Initiatives",
					"").
					AddOptionCheck(profile, acfn.InStrategyCommunity, true).
					AddOptionCheck(profile, acfn.InitiativeCommunityExists, true),
			).AddGroupOption(
			// Enables the user to create Initiative Communities
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateInitiativeCommunity"],
				"Initiative Communities",
				"").
				AddOptionCheck(profile, acfn.InStrategyCommunity, false).
				AddOptionCheck(profile, acfn.InCapabilityCommunity, true).
				AddOptionCheck(profile, acfn.CompanyVisionExists, true),
		).AddGroupOption(
			// Enables users from the strategy community to create Initiative Communities
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateInitiativeCommunity"],
				"Initiative Communities",
				"").
				AddOptionCheck(profile, acfn.InStrategyCommunity, true).
				AddOptionCheck(profile, acfn.CapabilityCommunityExists, true),
//***************************************************************************************************//
		).AddGroupOption(
			// Enables the user to create team values
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateValues"],
				"Competencies",
				"").
				AddOptionCheck(profile, acfn.InValuesCommunity, true),
		).AddGroupOption(
			// Enables the user to create company holidays
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateHolidays"],
				"Holidays",
				"").
				AddOptionCheck(profile, acfn.InHRCommunity, true),
		),
	).AddGroup(
		menu.NewGroupSpecification("Assign").
			AddGroupOption(
				// Enables a user to assign an Objective to an Objective Community
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["AssignCapabilityObjective"],
					"Assign Objective",
					"").
					AddOptionCheck(profile, acfn.MultipleCapabilityCommunitiesExists, true).
					AddOptionCheck(profile, acfn.ObjectivesExist, true).
					AddOptionCheck(profile, acfn.InStrategyCommunity, true),
			),
	).AddGroup(
		menu.NewGroupSpecification("Settings").
			AddGroupOption(
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["UserSettings"],
					"Update",
					""),
			),
	).AddGroup(
		menu.NewGroupSpecification("Reports").
			AddGroupOption(
				// Enables the user to see their last collaboration report
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["ViewCollaborationReport"],
					"Collaboration Report",
					"").
					AddOptionCheck(profile, acfn.CollaborationReportExists, true).
					AddOptionCheck(profile, acfn.InStrategyCommunity, true),
			).
			AddGroupOption(
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["StrategyPerformanceReport"],
					"Strategy Performance",
					"").
					AddOptionCheck(profile, acfn.InStrategyCommunity, true),
			).
			AddGroupOption(
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["IDOPerformanceReport"],
					"IDO Performance",
					"").
					AddChecks(p.Check("IDOsExistForMe")),
			),
	)
	return adm
}
