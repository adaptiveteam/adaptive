package adaptive_check_function_names

/* For every check function added it's name must be included in this list */
const (
	// Feedback tests
	FeedbackGivenThisQuarter = "FeedbackGivenThisQuarter"
	FeedbackForThePreviousQuarterExists = "FeedbackForThePreviousQuarterExists"
	InLastMonthOfQuarter     = "InLastMonthOfQuarter"

	CoacheesExist                             = "CoacheesExist"
	AdvocatesExist                            = "AdvocatesExist"

	// Due date tests
	IDOsDueWithinTheWeek    = "IDOsDueWithinTheWeek"
	IDOsDueWithinTheMonth   = "IDOsDueWithinTheMonth"
	IDOsDueWithinTheQuarter = "IDOsDueWithinTheQuarter"

	InitiativesDueWithinTheWeek    = "InitiativesDueWithinTheWeek"
	InitiativesDueWithinTheMonth   = "InitiativesDueWithinTheMonth"
	InitiativesDueWithinTheQuarter = "InitiativesDueWithinTheQuarter"

	ObjectivesDueWithinTheWeek    = "ObjectiveDueWithinTheWeek"
	ObjectivesDueWithinTheMonth   = "ObjectiveDueWithinTheMonth"
	ObjectivesDueWithinTheQuarter = "ObjectiveDueWithinTheQuarter"

	// Community membership tests
	InCapabilityCommunity = "InCapabilityCommunity"
	InValuesCommunity     = "InValuesCommunity"
	InHRCommunity         = "InHRCommunity"
	InStrategyCommunity   = "InStrategyCommunity"
	InInitiativeCommunity = "InInitiativeCommunity"

	// Component existence tests

	// Miscellaneous
	UserSettingsExist                = "UserSettingsExist"
	HolidaysExist                    = "HolidaysExist"
	CollaborationReportExists        = "CollaborationReportExists"
	UndeliveredEngagementsExistForMe = "UndeliveredEngagementsExistForMe"
	UndeliveredEngagementsOrPostponedEventsExistForMe = "UndeliveredEngagementsOrPostponedEventsExistForMe"
	CanBeNudgedForIDO                = "CanBeNudgedForIDO"

	// Strategy component existence tests independent of the user
	TeamValuesExist     = "TeamValuesExist"
	CompanyVisionExists = "CompanyVisionExists"
	ObjectivesExist     = "ObjectivesExist"
	InitiativesExist    = "InitiativesExist"

	// Strategy component existence tests for a given user
	IDOsExistForMe        = "IDOsExistForMe"
	ObjectivesExistForMe  = "ObjectivesExistForMe"
	InitiativesExistForMe = "InitiativesExistForMe"

	// Stale components exist for a specfc individual
	StaleInitiativesExistForMe = "StaleInitiativesExistForMe"
	StaleObjectivesExistForMe  = "StaleObjectivesExistForMe"

	// Community existence tests
	CapabilityCommunityExists           = "CapabilityCommunityExists"
	MultipleCapabilityCommunitiesExists = "MultipleCapabilityCommunitiesExists"
	InitiativeCommunityExists           = "InitiativeCommunityExists"
	MultipleInitiativeCommunitiesExists = "MultipleInitiativeCommunitiesExists"

	// State of community tests
	ObjectivesExistInMyCapabilityCommunities  = "ObjectivesExistInMyCapabilityCommunities"
	InitiativesExistInMyCapabilityCommunities = "InitiativesExistInMyCapabilityCommunities"
	InitiativesExistInMyInitiativeCommunities = "InitiativesExistInMyInitiativeCommunities"
)
