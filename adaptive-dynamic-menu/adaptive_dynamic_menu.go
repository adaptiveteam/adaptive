package adaptive_dynamic_menu

import (
	"github.com/adaptiveteam/adaptive/adaptive-checks"
	menu "github.com/adaptiveteam/adaptive/dynamic-menu"
)

func AdaptiveDynamicMenu(profile adaptive_checks.TypedProfile, bindings menu.FunctionBindings) (adm menu.DynamicMenuSpecification) {
	adm = menu.NewAdaptiveDynamicMenu().
	AddGroup(
		menu.NewGroupSpecification("Urgent Responsibilities").
			AddGroupOption(
				// Enables the user to create the company vision
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["CreateVision"],
					"Create Vision",
					"",
					!profile.CompanyVisionExists() &&
						profile.InStrategyCommunity(),
					),
			).AddGroupOption(
			// This fetches any undelivered engagements
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["FetchEngagementsForMe"],
				"What do I have right now?",
				"",
				profile.UndeliveredEngagementsOrPostponedEventsExistForMe(),
			),
		).AddGroupOption(
			// This fetches any IDO's not updated in the last 7 days
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["StaleIDOsExistForMe"],
				"Update IDO's",
				"",
				profile.StaleIDOsExistForMe(),
			),
		).AddGroupOption(
			// This fetches any Objectives not updated in the last 7 days
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["StaleObjectivesExistForMe"],
				"Update Objectives",
				"",
				profile.StaleObjectivesExistForMe(),
			),
		).AddGroupOption(
			// This fetches any Initiatives not updated in the last 7 days
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["StaleInitiativesExistForMe"],
				"Update Initiatives",
				"",
				profile.StaleInitiativesExistForMe(),
			),
		).AddGroupOption(
			// Enables the user to create an IDO
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateIDO"],
				"Create IDO",
				"",
				profile.CanBeNudgedForIDO() &&
				!profile.IDOsExistForMe() &&
				profile.TeamValuesExist(),
			),
		).AddGroupOption(
			// This enables the user to post feedback to another user
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ProvideFeedback"],
				"Provide Feedback",
				"",
				profile.InLastMonthOfQuarter() && 
				profile.TeamValuesExist(),
			),
		).AddGroupOption(
			// This enables the user to request feedback from another user
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["RequestFeedback"],
				"Request Feedback",
				"",
				profile.InLastMonthOfQuarter() &&
				profile.TeamValuesExist(),
			),
		),
	).AddGroup(
		menu.NewGroupSpecification("Responsibilities").
			AddGroupOption(
				// This fetches all IDO's
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["AllIDOsForMe"],
					"All IDO's",
					"",
					profile.IDOsExistForMe(),
				),
			).AddGroupOption(
			// This fetches all Objectives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["AllObjectivesForMe"],
				"All Objectives",
				"",
				profile.ObjectivesExistForMe(),
			),
		).AddGroupOption(
			// This fetches all Initiatives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["AllInitiativesForMe"],
				"All Initiatives",
				"",
				profile.InitiativesExistForMe(),
			),
		).AddGroupOption(
			// This enables the user to post feedback to another user
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ProvideFeedback"],
				"Provide Feedback",
				"",
				!profile.InLastMonthOfQuarter() &&
				profile.TeamValuesExist(),
			),
		).AddGroupOption(
			// This enables the user to request feedback from another user
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["RequestFeedback"],
				"Request Feedback",
				"",
				!profile.InLastMonthOfQuarter() &&
				profile.TeamValuesExist(),
			),
		),
	).AddGroup(
		menu.NewGroupSpecification("View").
			AddGroupOption(
				// Presents the company vision
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["ViewVision"],
					"Vision",
					"",
					profile.CompanyVisionExists() &&
					!profile.InStrategyCommunity(),
				),
			).AddGroupOption(
			// Presents the company vision
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewEditVision"],
				"Vision",
				"",
				profile.CompanyVisionExists() &&
				profile.InStrategyCommunity(),
			),
		).AddGroupOption(
			// Enables the user to see all of the objectives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewEditObjectives"],
				"Objectives",
				"",
				profile.ObjectivesExist() &&
				profile.InStrategyCommunity(),
			),
		).AddGroupOption(
			// Enables the user to see all of the objectives in their Capability & Initiative Communities
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewCommunityObjectives"],
				"Objectives",
				"",
				profile.ObjectivesExistInMyCapabilityCommunities() &&
				!profile.InStrategyCommunity(),
			),
		).AddGroupOption(
			// Enables the user to see all of the initiatives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewCommunityInitiatives"],
				"Initiatives",
				"",
				!profile.InStrategyCommunity() &&
				profile.InitiativesExistInMyInitiativeCommunities() &&
				!profile.InCapabilityCommunity(),
			),
		).AddGroupOption(
			// Enables the user to see and edit all of the Initiatives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewEditInitiatives"],
				"Initiatives",
				"",
				!profile.InStrategyCommunity() &&
				profile.InitiativesExistInMyCapabilityCommunities() &&
				profile.InCapabilityCommunity(),
			),
		).AddGroupOption(
			// Enables the user to see and edit all of the Initiatives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewEditInitiatives"],
				"Initiatives",
				"",
				profile.InStrategyCommunity() &&
				profile.InitiativesExist(),
			),
		).AddGroupOption(
			// Enables the user to see all of the team values
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewEditValues"],
				"Competencies",
				"",
				profile.TeamValuesExist() &&
				profile.InValuesCommunity(),
			),
		).AddGroupOption(
			// Enables the user to see all of the team values
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewValues"],
				"Competencies",
				"",
				profile.TeamValuesExist() &&
				!profile.InValuesCommunity(),
			),
		).AddGroupOption(
			// Enables the user to see all of the company holidays
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewHolidays"],
				"Holidays",
				"",
				profile.HolidaysExist() &&
				!profile.InHRCommunity(),
			),
		).AddGroupOption(
			// Enables the user to see all of the company holidays
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewEditHolidays"],
				"Holidays",
				"",
				profile.HolidaysExist() &&
				profile.InHRCommunity(),
			),
		).AddGroupOption(
			// Enables the user to see all of the scheduled events for the current quarter
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewScheduleCurrentQuarter"],
				"Current Quarter Events",
				"", 
				true,
			),
		).AddGroupOption(
			// Enables the user to see all of the scheduled events for the next quarter
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewScheduleNextQuarter"],
				"Next Quarter Events",
				"", 
				true,
			),
		).AddGroupOption(
			// Enables the user to see all of the IDOs they are coaching
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewCoachees"],
				"Coachee IDOs",
				"",
				profile.CoacheesExist(),
			),
		).AddGroupOption(
			// Enables the user to see all of the advocates for strategic elements aligned to you
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["ViewAdvocates"],
				"Advocates",
				"",
				profile.AdvocatesExist(),
			),
		),
	).AddGroup(
		menu.NewGroupSpecification("Create").
			AddGroupOption(
				// Enables a user to create an IDO
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["CreateIDO"],
					"IDO",
					"",
					profile.CanBeNudgedForIDO() &&
					profile.TeamValuesExist() &&
					profile.IDOsExistForMe(),
				),
			).AddGroupOption(
			// Enables a non user/initiative community user to create an IDO
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateIDO"],
				"IDO",
				"",
				!profile.CanBeNudgedForIDO() &&
				profile.TeamValuesExist(),),
		).AddGroupOption(
			// Enables the user to create objectives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateCapabilityObjectives"],
				"Objectives",
				"",
				profile.CompanyVisionExists() &&
				profile.InStrategyCommunity() &&
				profile.CapabilityCommunityExists(),
			),
		).AddGroupOption(
			// Enables the user to create Objective Communities
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateCapabilityCommunity"],
				"Objective Communities",
				"",
				profile.CompanyVisionExists() &&
				profile.InStrategyCommunity(),
			),
			//***************************************************************************************************//
			// The next four entries enable people who are either in a objective community to
			// create an initiative community or initiative if they are in the strategy community
			// or in an associated objective community.  The first set of two group options
			// enables people who are in a objective community to create an initiative community
			// (if there is a objective community) or an initiative
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
				"",
				!profile.InStrategyCommunity() &&
				profile.ObjectivesExistInMyCapabilityCommunities() &&
				profile.InInitiativeCommunity(),
			),
		).AddGroupOption(
			// Enables users from the strategy community to create Initiatives
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateInitiatives"],
				"Initiatives",
				"",
				profile.InStrategyCommunity() &&
				profile.InitiativeCommunityExists(),
			),
		).AddGroupOption(
			// Enables the user to create Initiative Communities
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateInitiativeCommunity"],
				"Initiative Communities",
				"",
				!profile.InStrategyCommunity() &&
				profile.InCapabilityCommunity() &&
				profile.CompanyVisionExists(),
			),
		).AddGroupOption(
			// Enables users from the strategy community to create Initiative Communities
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateInitiativeCommunity"],
				"Initiative Communities",
				"",
				profile.InStrategyCommunity() &&
				profile.CapabilityCommunityExists(),
			),
			//***************************************************************************************************//
		).AddGroupOption(
			// Enables the user to create team values
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateValues"],
				"Competencies",
				"",
				profile.InValuesCommunity(),
			),
		).AddGroupOption(
			// Enables the user to create company holidays
			menu.NewAdaptiveDynamicMenuSpecification(
				bindings["CreateHolidays"],
				"Holidays",
				"",
				profile.InHRCommunity(),
			),
		),
	).AddGroup(
		menu.NewGroupSpecification("Assign").
			AddGroupOption(
				// Enables a user to assign an Objective to an Objective Community
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["AssignCapabilityObjective"],
					"Assign Objective",
					"",
					profile.MultipleCapabilityCommunitiesExists() &&
					profile.ObjectivesExist() &&
					profile.InStrategyCommunity(),
				),
			),
	).AddGroup(
		menu.NewGroupSpecification("Settings").
			AddGroupOption(
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["UserSettings"],
					"Update",
					"",
					true,
				),
			),
	).AddGroup(
		menu.NewGroupSpecification("Reports").
			AddGroupOption(
				// Enables the user to see their last collaboration report
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["ViewCollaborationReport"],
					"Collaboration Report",
					"",
					profile.CollaborationReportExists(),
				),
				// profile.InStrategyCommunity(),),
			).
			AddGroupOption(
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["StrategyPerformanceReport"],
					"Strategy Performance",
					"",
					profile.InStrategyCommunity(),
				),
			).
			AddGroupOption(
				menu.NewAdaptiveDynamicMenuSpecification(
					bindings["IDOPerformanceReport"],
					"IDO Performance",
					"",
					profile.IDOsExistForMe(),
				),
			),
	)
	return adm
}
