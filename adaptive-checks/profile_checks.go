package adaptive_checks

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/core-utils-go"
	"golang.org/x/sync/errgroup"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/Merovius/go-misc/lazy"
)

type LazyBool = func ()bool

type TypedProfileConstructor = func (conn common.DynamoDBConnection, userID string, date business_time.Date) TypedProfile
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

var EvalProfile TypedProfileConstructor = func (conn common.DynamoDBConnection, userID string, date business_time.Date) TypedProfile {
	isUserInCommunity := isUserInCommunityCurry(conn, userID)

	env := readEnvironment()

	teamID := models.ParseTeamID(conn.PlatformID)
	return TypedProfile{
		FeedbackGivenThisQuarter: lazy.Bool(func() bool { return FeedbackGivenForTheQuarter(env, teamID, userID, date) }),
		FeedbackForThePreviousQuarterExists: lazy.Bool(func() bool {return FeedbackForThePreviousQuarterExists(env, teamID, userID, date)}),
		InLastMonthOfQuarter: lazy.Bool(func() bool {return date.GetMonth()%3 == 0 }),
		CoacheesExist: lazy.Bool(func() bool {return CoacheesExist(env, teamID, userID, date)}),
		AdvocatesExist: lazy.Bool(func() bool {return AdvocatesExist(env, teamID, userID, date)}),
		// IDOsDueWithinTheWeek: lazy.Bool(func() bool {return IDOsDueWithinTheWeek(env, teamID, userID, date)}),
		// IDOsDueWithinTheMonth: lazy.Bool(func() bool {return IDOsDueWithinTheMonth(env, teamID, userID, date)}),
		// IDOsDueWithinTheQuarter: lazy.Bool(func() bool {return IDOsDueWithinTheQuarter(env, teamID, userID, date)}),
		// InitiativesDueWithinTheWeek: lazy.Bool(func() bool {return InitiativesDueWithinTheWeek(env, teamID, userID, date)}),
		// InitiativesDueWithinTheMonth: lazy.Bool(func() bool {return InitiativesDueWithinTheMonth(env, teamID, userID, date)}),
		// InitiativesDueWithinTheQuarter: lazy.Bool(func() bool {return InitiativesDueWithinTheQuarter(env, teamID, userID, date)}),
		// ObjectivesDueWithinTheWeek: lazy.Bool(func() bool {return ObjectivesDueWithinTheWeek(env, teamID, userID, date)}),
		// ObjectivesDueWithinTheMonth: lazy.Bool(func() bool {return ObjectivesDueWithinTheMonth(env, teamID, userID, date)}),
		// ObjectivesDueWithinTheQuarter: lazy.Bool(func() bool {return ObjectivesDueWithinTheQuarter(env, teamID, userID, date)}),
		InCapabilityCommunity: lazy.Bool(func() bool {return isUserInCommunity(community.Capability)}),
		InValuesCommunity: lazy.Bool(func() bool {return isUserInCommunity(community.Competency)}),
		InHRCommunity: lazy.Bool(func() bool {return isUserInCommunity(community.HR)}),
		InStrategyCommunity: lazy.Bool(func() bool {return isUserInCommunity(community.Strategy) }),
		InInitiativeCommunity: lazy.Bool(func() bool {return InitiativeCommunityExistsForMe(env, teamID, userID, date)}),
		UserSettingsExist: lazy.Bool(func() bool {return true}),
		HolidaysExist: lazy.Bool(func() bool {return HolidaysExist(env, teamID, userID, date)}),
		CollaborationReportExists: lazy.Bool(func() bool {return ReportExists(env, teamID, userID, date)}),
		UndeliveredEngagementsExistForMe: lazy.Bool(func() bool {return UndeliveredEngagementsExistForMe(env, teamID, userID, date)}),
		UndeliveredEngagementsOrPostponedEventsExistForMe: lazy.Bool(func() bool {return UndeliveredEngagementsOrPostponedEventsExistForMe(env, teamID, userID, date)}),
		CanBeNudgedForIDO: lazy.Bool(func() bool {return CanBeNudgedForIDOCreation(env, teamID, userID, date)}),
		TeamValuesExist: lazy.Bool(func() bool {return TeamValuesExist(env, teamID, userID, date)}),
		CompanyVisionExists: lazy.Bool(func() bool {return CompanyVisionExists(env, teamID, userID, date)}),
		ObjectivesExist: lazy.Bool(func() bool {return ObjectivesExist(env, teamID, userID, date)}),
		InitiativesExist: lazy.Bool(func() bool {return InitiativesExistInMyCapabilityCommunities(env, teamID, userID, date)}),
		IDOsExistForMe: lazy.Bool(func() bool {return IDOsExistForMe(env, teamID, userID, date)}),
		ObjectivesExistForMe: lazy.Bool(func() bool {return ObjectivesExistForMe(env, teamID, userID, date)}),
		InitiativesExistForMe: lazy.Bool(func() bool {return InitiativesExistForMe(env, teamID, userID, date)}),
		StaleIDOsExistForMe: lazy.Bool(func() bool {return StaleIDOsExist(env, teamID, userID, date)}),
		StaleInitiativesExistForMe: lazy.Bool(func() bool {return StaleInitiativesExistForMe(env, teamID, userID, date)}),
		StaleObjectivesExistForMe: lazy.Bool(func() bool {return StaleObjectivesExistForMe(env, teamID, userID, date)}),
		CapabilityCommunityExists: lazy.Bool(func() bool {return CapabilityCommunityExists(env, teamID, userID, date)}),
		MultipleCapabilityCommunitiesExists: lazy.Bool(func() bool {return MultipleCapabilityCommunitiesExists(env, teamID, userID, date)}),
		InitiativeCommunityExists: lazy.Bool(func() bool {return InitiativeCommunityExists(env, teamID, userID, date)}),
		MultipleInitiativeCommunitiesExists: lazy.Bool(func() bool {return false}),
		ObjectivesExistInMyCapabilityCommunities: lazy.Bool(func() bool {return ObjectivesExistInMyCapabilityCommunities(env, teamID, userID, date)}),
		InitiativesExistInMyCapabilityCommunities: lazy.Bool(func() bool {return InitiativesExistInMyCapabilityCommunities(env, teamID, userID, date)}),
		InitiativesExistInMyInitiativeCommunities: lazy.Bool(func() bool {return InitiativesExistInMyInitiativeCommunities(env, teamID, userID, date)}),
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

func ConstLazyBool(value bool) func () bool {
	return func () bool { return value }
}

func ConstLazyInt(value int) func () int {
	return func () int { return value }
}
