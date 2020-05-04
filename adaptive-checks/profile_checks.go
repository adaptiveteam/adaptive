package adaptive_checks

import (
	"github.com/adaptiveteam/adaptive/core-utils-go"
	"golang.org/x/sync/errgroup"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
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

func LazyUserCommunities(f func ()[]adaptiveCommunityUser.AdaptiveCommunityUser) func ()[]adaptiveCommunityUser.AdaptiveCommunityUser {
	lazyI := lazy.Interface(func () interface{} {
		return f()
	})
	return func ()[]adaptiveCommunityUser.AdaptiveCommunityUser {
		i := lazyI()
		return i.([]adaptiveCommunityUser.AdaptiveCommunityUser)
	}
}

func isUserInCommunityCurry(conn common.DynamoDBConnection, userID string) func (communityID community.AdaptiveCommunity) bool {
	userCommunities := LazyUserCommunities(func ()[]adaptiveCommunityUser.AdaptiveCommunityUser {
		return adaptiveCommunityUser.ReadByUserIDUnsafe(userID)(conn)
	})
	return func (communityID community.AdaptiveCommunity) bool {
		for _, uc := range userCommunities() {
			if uc.CommunityID == string(communityID) {
				return true
			}
		}
		return false
	}
}

func EvalProfile(conn common.DynamoDBConnection, userID string, date business_time.Date) TypedProfile {
	isUserInCommunity := isUserInCommunityCurry(conn, userID)

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
		InCapabilityCommunity: lazy.Bool(func() bool {return isUserInCommunity(community.Capability)}),
		InValuesCommunity: lazy.Bool(func() bool {return isUserInCommunity(community.Competency)}),
		InHRCommunity: lazy.Bool(func() bool {return isUserInCommunity(community.HR)}),
		InStrategyCommunity: lazy.Bool(func() bool {return isUserInCommunity(community.Strategy) }),
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

// EagerLoad loads the given list of lazy bools in parallel
func EagerLoad(lazyBools ... LazyBool) error {
	var errGroup errgroup.Group
	for _, l := range lazyBools {
		errGroup.Go(func() (err1 error) {
			defer core_utils_go.RecoverToErrorVar("EagerLoad", &err1)
			l()
			return 
		})
	}

	return errGroup.Wait()
}

// LoadAll forces all lazy evaluations
func (p *TypedProfile) LoadAll() error {
	return EagerLoad(
		p.FeedbackGivenThisQuarter,
		p.FeedbackForThePreviousQuarterExists,
		p.InLastMonthOfQuarter,
		p.CoacheesExist,
		p.AdvocatesExist,
		p.InCapabilityCommunity,
		p.InValuesCommunity,
		p.InHRCommunity,
		p.InStrategyCommunity,
		p.InInitiativeCommunity,
		p.UserSettingsExist,
		p.HolidaysExist,
		p.CollaborationReportExists,
		p.UndeliveredEngagementsExistForMe,
		p.UndeliveredEngagementsOrPostponedEventsExistForMe,
		p.CanBeNudgedForIDO,
		p.TeamValuesExist,
		p.CompanyVisionExists,
		p.ObjectivesExist,
		p.InitiativesExist,
		p.IDOsExistForMe,
		p.ObjectivesExistForMe,
		p.InitiativesExistForMe,
		p.StaleIDOsExistForMe,
		p.StaleInitiativesExistForMe,
		p.StaleObjectivesExistForMe,
		p.CapabilityCommunityExists,
		p.MultipleCapabilityCommunitiesExists,
		p.InitiativeCommunityExists,
		p.MultipleInitiativeCommunitiesExists,
		p.ObjectivesExistInMyCapabilityCommunities,
		p.InitiativesExistInMyCapabilityCommunities,
		p.InitiativesExistInMyInitiativeCommunities,
	)
}