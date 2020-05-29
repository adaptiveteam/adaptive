package adaptive_dynamic_menu

import (
	adaptive_checks "github.com/adaptiveteam/adaptive/adaptive-checks"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	menu "github.com/adaptiveteam/adaptive/dynamic-menu"
)

func AdaptiveDynamicMenu(profile adaptive_checks.TypedProfile, bindings menu.FunctionBindings) (adm menu.DynamicMenuSpecification) {
	adm = menu.NewAdaptiveDynamicMenu(
		UrgentResponsibilitiesGroup(profile, bindings),
		ResponsibilitiesGroup(profile, bindings),
		ViewGroup(profile, bindings),
		CreateGroup(profile, bindings),
		AssignGroup(profile, bindings),
		SettingsGroup(profile, bindings),
		ReportsGroup(profile, bindings),
	)
	return adm
}

func UrgentResponsibilitiesGroup(profile adaptive_checks.TypedProfile, bindings menu.FunctionBindings) menu.GroupSpecification {
	return menu.NewGroupSpecification(
		"Urgent Responsibilities",
		// Enables the user to create the company vision
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["CreateVision"],
			"Create Vision",
			"",
			!profile.CompanyVisionExists() &&
				profile.InStrategyCommunity(),
		),
		// This fetches any undelivered engagements
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["FetchEngagementsForMe"],
			"What do I have right now?",
			"",
			profile.UndeliveredEngagementsOrPostponedEventsExistForMe(),
		),
		// This fetches any IDO's not updated in the last 7 days
		menu.NewAdaptiveDynamicMenuSpecification(
			user.StaleIDOsForMe,
			"Update IDO's",
			"",
			profile.StaleIDOsExistForMe(),
		),
		// This fetches any Objectives not updated in the last 7 days
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["StaleObjectivesExistForMe"],
			"Update Objectives",
			"",
			profile.StaleObjectivesExistForMe(),
		),
		// This fetches any Initiatives not updated in the last 7 days
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["StaleInitiativesExistForMe"],
			"Update Initiatives",
			"",
			profile.StaleInitiativesExistForMe(),
		),
		// Enables the user to create an IDO
		menu.NewAdaptiveDynamicMenuSpecification(
			objectives.CreateIDONow,
			"Create IDO",
			"",
			profile.CanBeNudgedForIDO() &&
				!profile.IDOsExistForMe() &&
				profile.TeamValuesExist(),
		),
		// This enables the user to post feedback to another user
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["ProvideFeedback"],
			"Provide Feedback",
			"",
			profile.InLastMonthOfQuarter() &&
				profile.TeamValuesExist(),
		),
		// This enables the user to request feedback from another user
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["RequestFeedback"],
			"Request Feedback",
			"",
			profile.InLastMonthOfQuarter() &&
				profile.TeamValuesExist(),
		),
	)
}

func ResponsibilitiesGroup(profile adaptive_checks.TypedProfile, bindings menu.FunctionBindings) menu.GroupSpecification {
	return menu.NewGroupSpecification(
		"Responsibilities",
		// This fetches all IDO's
		menu.NewAdaptiveDynamicMenuSpecification(
			user.ViewObjectives,
			"All IDO's",
			"",
			profile.IDOsExistForMe(),
		),
		// This fetches all Objectives
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["AllObjectivesForMe"],
			"All Objectives",
			"",
			profile.ObjectivesExistForMe(),
		),
		// This fetches all Initiatives
		menu.NewAdaptiveDynamicMenuSpecification(
			strategy.ViewAdvocacyInitiatives,
			"All Initiatives",
			"",
			profile.InitiativesExistForMe(),
		),
		// This enables the user to post feedback to another user
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["ProvideFeedback"],
			"Provide Feedback",
			"",
			!profile.InLastMonthOfQuarter() &&
				profile.TeamValuesExist(),
		),
		// This enables the user to request feedback from another user
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["RequestFeedback"],
			"Request Feedback",
			"",
			!profile.InLastMonthOfQuarter() &&
				profile.TeamValuesExist(),
		),
	)
}

func ViewGroup(profile adaptive_checks.TypedProfile, bindings menu.FunctionBindings) menu.GroupSpecification {
	return menu.NewGroupSpecification(
		"View",
		// Presents the company vision
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["ViewVision"],
			"Vision",
			"",
			profile.CompanyVisionExists() &&
				!profile.InStrategyCommunity(),
		),
		// Presents the company vision
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["ViewEditVision"],
			"Vision",
			"",
			profile.CompanyVisionExists() &&
				profile.InStrategyCommunity(),
		),
		// Enables the user to see all of the objectives
		menu.NewAdaptiveDynamicMenuSpecification(
			strategy.ViewStrategyObjectives,
			"Objectives",
			"",
			profile.ObjectivesExist() &&
				profile.InStrategyCommunity(),
		),
		// Enables the user to see all of the objectives in their Capability & Initiative Communities
		menu.NewAdaptiveDynamicMenuSpecification(
			strategy.ViewCapabilityCommunityObjectives,
			"Objectives",
			"",
			profile.ObjectivesExistInMyCapabilityCommunities() &&
				!profile.InStrategyCommunity(),
		),
		// Enables the user to see all of the initiatives
		menu.NewAdaptiveDynamicMenuSpecification(
			strategy.ViewInitiativeCommunityInitiatives,
			"Initiatives",
			"",
			!profile.InStrategyCommunity() &&
				profile.InitiativesExistInMyInitiativeCommunities() &&
				!profile.InCapabilityCommunity(),
		),
		// Enables the user to see and edit all of the Initiatives
		menu.NewAdaptiveDynamicMenuSpecification(
			strategy.ViewCapabilityCommunityInitiatives,
			"Initiatives",
			"",
			(!profile.InStrategyCommunity() &&
				profile.InitiativesExistInMyCapabilityCommunities() &&
				profile.InCapabilityCommunity()) ||
			(profile.InStrategyCommunity() &&
				profile.InitiativesExist()),
		),
		// Enables the user to see all of the team values
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["ViewEditValues"],
			"Competencies",
			"",
			profile.TeamValuesExist() &&
				profile.InValuesCommunity(),
		),
		// Enables the user to see all of the team values
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["ViewValues"],
			"Competencies",
			"",
			profile.TeamValuesExist() &&
				!profile.InValuesCommunity(),
		),
		// Enables the user to see all of the company holidays
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["ViewHolidays"],
			"Holidays",
			"",
			profile.HolidaysExist() &&
				!profile.InHRCommunity(),
		),
		// Enables the user to see all of the company holidays
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["ViewEditHolidays"],
			"Holidays",
			"",
			profile.HolidaysExist() &&
				profile.InHRCommunity(),
		),
		// Enables the user to see all of the scheduled events for the current quarter
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["ViewScheduleCurrentQuarter"],
			"Current Quarter Events",
			"",
			true,
		),
		// Enables the user to see all of the scheduled events for the next quarter
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["ViewScheduleNextQuarter"],
			"Next Quarter Events",
			"",
			true,
		),
		// Enables the user to see all of the IDOs they are coaching
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["ViewCoachees"],
			"Coachee IDOs",
			"",
			profile.CoacheesExist(),
		),
		// Enables the user to see all of the advocates for strategic elements aligned to you
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["ViewAdvocates"],
			"Advocates",
			"",
			profile.AdvocatesExist(),
		),
	)
}

func CreateGroup(profile adaptive_checks.TypedProfile, bindings menu.FunctionBindings) menu.GroupSpecification {
	return menu.NewGroupSpecification(
		"Create",
		// Enables a user to create an IDO
		menu.NewAdaptiveDynamicMenuSpecification(
			objectives.CreateIDONow,
			"IDO",
			"",
			(profile.CanBeNudgedForIDO() &&
				profile.TeamValuesExist() &&
				profile.IDOsExistForMe()) ||
			(!profile.CanBeNudgedForIDO() &&
				profile.TeamValuesExist()),
		),
		// Enables the user to create objectives
		menu.NewAdaptiveDynamicMenuSpecification(
			strategy.CreateStrategyObjective,
			"Objectives",
			"",
			profile.CompanyVisionExists() &&
				profile.InStrategyCommunity() &&
				profile.CapabilityCommunityExists(),
		),
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
		// Enables the user to create Initiatives
		menu.NewAdaptiveDynamicMenuSpecification(
			strategy.CreateInitiative,
			"Initiatives",
			"",
			(!profile.InStrategyCommunity() &&
				profile.ObjectivesExistInMyCapabilityCommunities() &&
				profile.InInitiativeCommunity()) || 
				(profile.InStrategyCommunity() &&
				profile.InitiativeCommunityExists()),
		),
		// Enables the user to create Initiative Communities
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["CreateInitiativeCommunity"],
			"Initiative Communities",
			"",
			(!profile.InStrategyCommunity() &&
				profile.InCapabilityCommunity() &&
				profile.CompanyVisionExists()) || 
			(profile.InStrategyCommunity() &&
				profile.CapabilityCommunityExists()),
		),
		//***************************************************************************************************//
		// Enables the user to create team values
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["CreateValues"],
			"Competencies",
			"",
			profile.InValuesCommunity(),
		),
		// Enables the user to create company holidays
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["CreateHolidays"],
			"Holidays",
			"",
			profile.InHRCommunity(),
		),
	)
}

func AssignGroup(profile adaptive_checks.TypedProfile, bindings menu.FunctionBindings) menu.GroupSpecification {
	return menu.NewGroupSpecification(
		"Assign",
		// Enables a user to assign an Objective to an Objective Community
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["AssignCapabilityObjective"],
			"Assign Objective",
			"",
			profile.MultipleCapabilityCommunitiesExists() &&
				profile.ObjectivesExist() &&
				profile.InStrategyCommunity(),
		),
	)
}

func SettingsGroup(profile adaptive_checks.TypedProfile, bindings menu.FunctionBindings) menu.GroupSpecification {
	return menu.NewGroupSpecification(
		"Settings",
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["UserSettings"],
			"Update",
			"",
			true,
		),
	)
}

func ReportsGroup(profile adaptive_checks.TypedProfile, bindings menu.FunctionBindings) menu.GroupSpecification {
	return menu.NewGroupSpecification("Reports",
		// Enables the user to see their last collaboration report
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["ViewCollaborationReport"],
			"Collaboration Report",
			"",
			profile.CollaborationReportExists(),
		),
		// profile.InStrategyCommunity(),),
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["StrategyPerformanceReport"],
			"Strategy Performance",
			"",
			profile.InStrategyCommunity(),
		),
		menu.NewAdaptiveDynamicMenuSpecification(
			bindings["IDOPerformanceReport"],
			"IDO Performance",
			"",
			profile.IDOsExistForMe(),
		),
	)
}
