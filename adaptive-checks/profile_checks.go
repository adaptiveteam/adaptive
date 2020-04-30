package adaptive_checks

import (
	"github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/Merovius/go-misc/lazy"
)

type LazyBool = func ()bool

type TypedProfile struct {
	FeedbackGivenThisQuarter LazyBool
	FeedbackForThePreviousQuarterExists LazyBool
	InLastMonthOfQuarter LazyBool
	CoacheesExist LazyBool
	AdvocatesExist LazyBool
	IDOsDueWithinTheWeek LazyBool
	IDOsDueWithinTheMonth LazyBool
	IDOsDueWithinTheQuarter LazyBool
	InitiativesDueWithinTheWeek LazyBool
	InitiativesDueWithinTheMonth LazyBool
	InitiativesDueWithinTheQuarter LazyBool
	ObjectivesDueWithinTheWeek LazyBool
	ObjectivesDueWithinTheMonth LazyBool
	ObjectivesDueWithinTheQuarter LazyBool
	InCapabilityCommunity LazyBool
	InValuesCommunity LazyBool
	InHRCommunity LazyBool
	InStrategyCommunity LazyBool
	InInitiativeCommunity LazyBool
	UserSettingsExist LazyBool
	HolidaysExist LazyBool
	CollaborationReportExists LazyBool
	UndeliveredEngagementsExistForMe LazyBool
	UndeliveredEngagementsOrPostponedEventsExistForMe LazyBool
	CanBeNudgedForIDO LazyBool
	TeamValuesExist LazyBool
	CompanyVisionExists LazyBool
	ObjectivesExist LazyBool
	InitiativesExist LazyBool
	IDOsExistForMe LazyBool
	ObjectivesExistForMe LazyBool
	InitiativesExistForMe LazyBool
	StaleIDOsExistForMe LazyBool
	StaleInitiativesExistForMe LazyBool
	StaleObjectivesExistForMe LazyBool
	CapabilityCommunityExists LazyBool
	MultipleCapabilityCommunitiesExists LazyBool
	InitiativeCommunityExists LazyBool
	MultipleInitiativeCommunitiesExists LazyBool
	ObjectivesExistInMyCapabilityCommunities LazyBool
	InitiativesExistInMyCapabilityCommunities LazyBool
	InitiativesExistInMyInitiativeCommunities LazyBool
}

func EvalProfile(conn common.DynamoDBConnection, userID string, date business_time.Date) TypedProfile {
	return TypedProfile{
		FeedbackGivenThisQuarter: lazy.Bool(func() bool { return FeedbackGivenForTheQuarter(userID, date) }),
		FeedbackForThePreviousQuarterExists: lazy.Bool(func() bool {return FeedbackForThePreviousQuarterExists(userID, date)}),
		InLastMonthOfQuarter: lazy.Bool(func() bool {return date.GetMonth()%3 == 0 }),
		CoacheesExist: lazy.Bool(func() bool {return CoacheesExist(userID, date)}),
		AdvocatesExist: lazy.Bool(func() bool {return AdvocatesExist(userID, date)}),
		// IDOsDueWithinTheWeek: lazy.Bool(func() bool {return IDOsDueWithinTheWeek(userID, date)}),
		// IDOsDueWithinTheMonth: lazy.Bool(func() bool {return IDOsDueWithinTheMonth(userID, date)}),
		// IDOsDueWithinTheQuarter: lazy.Bool(func() bool {return IDOsDueWithinTheQuarter(userID, date)}),
		// InitiativesDueWithinTheWeek: lazy.Bool(func() bool {return InitiativesDueWithinTheWeek(userID, date)}),
		// InitiativesDueWithinTheMonth: lazy.Bool(func() bool {return InitiativesDueWithinTheMonth(userID, date)}),
		// InitiativesDueWithinTheQuarter: lazy.Bool(func() bool {return InitiativesDueWithinTheQuarter(userID, date)}),
		// ObjectivesDueWithinTheWeek: lazy.Bool(func() bool {return ObjectivesDueWithinTheWeek(userID, date)}),
		// ObjectivesDueWithinTheMonth: lazy.Bool(func() bool {return ObjectivesDueWithinTheMonth(userID, date)}),
		// ObjectivesDueWithinTheQuarter: lazy.Bool(func() bool {return ObjectivesDueWithinTheQuarter(userID, date)}),
		InCapabilityCommunity: lazy.Bool(func() bool {return InCapabilityCommunity(userID, date)}),
		InValuesCommunity: lazy.Bool(func() bool {return InCompetenciesCommunity(userID, date)}),
		InHRCommunity: lazy.Bool(func() bool {return InHRCommunity(userID, date)}),
		InStrategyCommunity: lazy.Bool(func() bool {return InStrategyCommunity(userID, date)}),
		InInitiativeCommunity: lazy.Bool(func() bool {return InitiativeCommunityExistsForMe(userID, date)}),
		UserSettingsExist: lazy.Bool(func() bool {return true}),
		HolidaysExist: lazy.Bool(func() bool {return HolidaysExist(userID, date)}),
		CollaborationReportExists: lazy.Bool(func() bool {return ReportExists(userID, date)}),
		UndeliveredEngagementsExistForMe: lazy.Bool(func() bool {return UndeliveredEngagementsExistForMe(userID, date)}),
		UndeliveredEngagementsOrPostponedEventsExistForMe: lazy.Bool(func() bool {return UndeliveredEngagementsOrPostponedEventsExistForMe(userID, date)}),
		CanBeNudgedForIDO: lazy.Bool(func() bool {return CanBeNudgedForIDOCreation(userID, date)}),
		TeamValuesExist: lazy.Bool(func() bool {return TeamValuesExist(userID, date)}),
		CompanyVisionExists: lazy.Bool(func() bool {return CompanyVisionExists(userID, date)}),
		ObjectivesExist: lazy.Bool(func() bool {return ObjectivesExist(userID, date)}),
		InitiativesExist: lazy.Bool(func() bool {return InitiativesExistInMyCapabilityCommunities(userID, date)}),
		IDOsExistForMe: lazy.Bool(func() bool {return IDOsExistForMe(userID, date)}),
		ObjectivesExistForMe: lazy.Bool(func() bool {return ObjectivesExistForMe(userID, date)}),
		InitiativesExistForMe: lazy.Bool(func() bool {return InitiativesExistForMe(userID, date)}),
		StaleIDOsExistForMe: lazy.Bool(func() bool {return StaleIDOsExist(userID, date)}),
		StaleInitiativesExistForMe: lazy.Bool(func() bool {return StaleInitiativesExistForMe(userID, date)}),
		StaleObjectivesExistForMe: lazy.Bool(func() bool {return StaleObjectivesExistForMe(userID, date)}),
		CapabilityCommunityExists: lazy.Bool(func() bool {return CapabilityCommunityExists(userID, date)}),
		MultipleCapabilityCommunitiesExists: lazy.Bool(func() bool {return MultipleCapabilityCommunitiesExists(userID, date)}),
		InitiativeCommunityExists: lazy.Bool(func() bool {return InitiativeCommunityExists(userID, date)}),
		MultipleInitiativeCommunitiesExists: lazy.Bool(func() bool {return false}),
		ObjectivesExistInMyCapabilityCommunities: lazy.Bool(func() bool {return ObjectivesExistInMyCapabilityCommunities(userID, date)}),
		InitiativesExistInMyCapabilityCommunities: lazy.Bool(func() bool {return InitiativesExistInMyCapabilityCommunities(userID, date)}),
		InitiativesExistInMyInitiativeCommunities: lazy.Bool(func() bool {return InitiativesExistInMyInitiativeCommunities(userID, date)}),
	}
}