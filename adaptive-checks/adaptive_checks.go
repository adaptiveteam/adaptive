package adaptive_checks

import (
	"log"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	eholidays "github.com/adaptiveteam/adaptive/adaptive-engagements/holidays"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/values"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	// "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/postponedEvent"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
)

/* IDO Checks */

const logEnabled = false

// IDOsExistForMe Are there any IDO's that exist for the user?
func IDOsExistForMe(env environment, teamID models.TeamID, userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("IDOsExistForMe")
	if logEnabled {
		log.Println("Checking IDOsExistForMe")
	}
	objs := objectives.AllUserObjectives(userID, env.userObjectivesTable, 
		env.userObjectivesUserIndex,
		models.IndividualDevelopmentObjective, 0)
	res = len(objs) > 0
	if logEnabled {
		log.Printf("IDOsExistForMe(%s, _): %v\n", userID, res)
	}
	return
}

// IDOsDueInAWeek Are there any open IDO's that exist for the user that are due in exactly 7 days
func IDOsDueInAWeek(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("IDOsDueInAWeek")
	op := objectives.IDOsDueInAWeek(userID, date, 
		env.userObjectivesTable, env.userObjectivesUserIndex)
	res = len(op) > 0
	if logEnabled {
		log.Printf("IDOsDueInAWeek(%s, %v): %v\n", userID, date, res)
	}
	return
}

// IDOsDueInAMonth Are there any open IDO's that exist for the user that are due in exactly in 30 days
func IDOsDueInAMonth(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("IDOsDueInAMonth")
	op := objectives.IDOsDueInAMonth(userID, date, 
		env.userObjectivesTable, env.userObjectivesUserIndex)
	res = len(op) > 0
	if logEnabled {
		log.Printf("IDOsDueInAMonth(%s, %v): %v\n", userID, date, res)
	}
	return
}

// IDOsDueInAQuarter Are there any open IDO's that exist for the user that are due in exactly 90 days
func IDOsDueInAQuarter(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("IDOsDueInAQuarter")
	op := objectives.IDOsDueInAQuarter(userID, date, 
		env.userObjectivesTable, env.userObjectivesUserIndex)
	res = len(op) > 0
	if logEnabled {
		log.Printf("IDOsDueInAQuarter(%s, %v): %v\n", userID, date, res)
	}
	return
}

// StaleIDOs checks if there are any stale IDO's for a user
func StaleIDOs(env environment, teamID models.TeamID, userID string, date business_time.Date) []models.UserObjective {
	defer RecoverToLog("StaleIDOs")
	return objectives.UserIDOsWithNoProgressInAWeek(userID, date, 
		env.userObjectivesTable, env.userObjectivesUserIndex,
		env.userObjectivesProgressTable)
}

// StaleIDOsExist checks that there are IDOs without recent progress
func StaleIDOsExist(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("StaleIDOsExist")
	staleIDOs := StaleIDOs(env, teamID, userID, date)
	res = len(staleIDOs) > 0
	if logEnabled {
		log.Printf("StaleIDOsExist(%s, %v): %v\n", userID, date, res)
	}
	return
}

/* Vision Checks */

// CompanyVisionExists Does the company vision exist?
func CompanyVisionExists(env environment, teamID models.TeamID, userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("CompanyVisionExists")
	return strategy.StrategyVision(teamID, env.visionTable) != nil
}

// InStrategyCommunity Is the user in the strategy community
func InStrategyCommunity(env environment, teamID models.TeamID, userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InStrategyCommunity")
	return community.IsUserInCommunity(userID, 
		env.communityUsersTable, env.communityUsersUserCommunityIndex, community.Strategy)
}

/* Objective Checks */

// ObjectivesExistForMe Is the user the advocate for any objectives?
func ObjectivesExistForMe(env environment, teamID models.TeamID, userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("ObjectivesExistForMe")
	objs := strategy.UserAdvocacyObjectives(userID, 
		env.userObjectivesTable, env.userObjectivesTypeIndex, 0)
	if logEnabled {
		log.Println("Checked ObjectivesExistForMe: ", objs)
	}
	return len(objs) > 0
}

// ObjectivesExist returns all the objectives for the organization
func ObjectivesExist(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("ObjectivesExist")
	if logEnabled {
		log.Printf("Checking ObjectivesExist for userID=%s, date=%v\n", userID, date)
	}
	platformID := teamID.ToPlatformID()
	conn := daosCommon.CreateConnectionFromEnv(platformID)
	pager := strategy.SelectFromStrategyObjectiveJoinCommunityWhereUserIDOrInStrategyCommunityStream(userID)(conn)
	var err error
	res, err = pager.NonEmpty()
	if err != nil {
		log.Printf("ERROR ObjectivesExist: %+v\n", err)
	}
	if logEnabled {
		log.Println("Checked ObjectivesExist: ", res)
	}
	return
}

// Is the user the advocate for any objectives?
// that have no been updated within the last month?
func StaleObjectivesExistForMe(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("StaleObjectivesExistForMe")
	stratObjs := strategy.UserCapabilityObjectivesWithNoProgressInAMonth(userID, date,
		env.userObjectivesTable, env.userObjectivesUserIndex, env.userObjectivesProgressTable, 0)
	if logEnabled {
		log.Println("Checked StaleObjectivesExistForMe: ", len(stratObjs))
	}
	return len(stratObjs) > 0
}

// ObjectivesExistInMyCapabilityCommunities checks
//   if the user belong to any capability communities that have
// Capability Objectives allocated to them?
func ObjectivesExistInMyCapabilityCommunities(env environment, 
	teamID models.TeamID, userID string, date business_time.Date,
) (res bool) {
	defer RecoverToLog("ObjectivesExistInMyCapabilityCommunities")
	if logEnabled {
		log.Printf("Checking ObjectivesExistInMyCapabilityCommunities for userID=%s, date=%v\n", userID, date)
	}
	conn := env.connGen.ForPlatformID(teamID.ToPlatformID())
	objs := strategy.UserCommunityObjectives(userID,
		env.strategyObjectivesTableName, env.strategyObjectivesPlatformIndex,
		env.userObjectivesTable,
		env.communityUsersTable, env.communityUsersUserIndex, conn)
	if logEnabled {
		log.Printf("Checked ObjectivesExistInMyCapabilityCommunities: %d\n", len(objs))
	}
	return len(objs) > 0
}

// CapabilityObjectivesDueInAWeek Are there any open Objectives that exist for the user that are due in exactly 7 days
func CapabilityObjectivesDueInAWeek(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("CapabilityObjectivesDueInAWeek")
	op := strategy.CapabilityObjectivesDueInAWeek(userID, date, env.userObjectivesTable, env.userObjectivesUserIndex)
	if logEnabled {
		log.Println("Checked CapabilityObjectivesDueInAWeek: ", len(op))
	}
	return len(op) > 0
}

// CapabilityObjectivesDueInAMonth Are there any open Objectives that exist for the user that are due in exactly 30 days
func CapabilityObjectivesDueInAMonth(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("CapabilityObjectivesDueInAMonth")
	op := strategy.CapabilityObjectivesDueInAMonth(userID, date, env.userObjectivesTable, env.userObjectivesUserIndex)
	if logEnabled {
		log.Println("Checked CapabilityObjectivesDueInAMonth: ", len(op))
	}
	return len(op) > 0
}

// CapabilityObjectivesDueInAQuarter Are there any open Objectives that exist for the user that are due in exactly 90 days
func CapabilityObjectivesDueInAQuarter(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("CapabilityObjectivesDueInAQuarter")
	op := strategy.CapabilityObjectivesDueInAQuarter(userID, date, env.userObjectivesTable, env.userObjectivesUserIndex)
	if logEnabled {
		log.Println("Checked CapabilityObjectivesDueInAQuarter: ", len(op))
	}
	return len(op) > 0
}

/* Capabilitity Community Checks */

// InCapabilityCommunity Is the user in any Objective Community?
func InCapabilityCommunity(env environment, teamID models.TeamID, userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InCapabilityCommunity")
	capComms, _ := strategy.UserCapabilityInitiativeCommunities(userID, 
		env.communityUsersTable, env.communityUsersUserIndex)
	if logEnabled {
		log.Println("Checked InCapabilityCommunity: ", len(capComms))
	}
	return len(capComms) > 0
}

// CapabilityCommunityExists Does there exist a capabilility community?
func CapabilityCommunityExists(env environment, teamID models.TeamID, userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("CapabilityCommunityExists")
	capComms := strategy.AllCapabilityCommunities(teamID, env.capabilityCommunitiesTable,
		env.capabilityCommunitiesPlatformIndex, env.strategyCommunitiesTable)
	if logEnabled {
		log.Println("Checked CapabilityCommunityExists: ", len(capComms))
	}
	return len(capComms) > 0
}

// MultipleCapabilityCommunitiesExists Is there more than one objective community?
func MultipleCapabilityCommunitiesExists(env environment, teamID models.TeamID, userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("MultipleCapabilityCommunitiesExists")
	capComms := strategy.AllCapabilityCommunities(teamID, env.capabilityCommunitiesTable,
		env.capabilityCommunitiesPlatformIndex, env.strategyCommunitiesTable)
	if logEnabled {
		log.Println("Checked MultipleCapabilityCommunitiesExists: ", len(capComms))
	}
	return len(capComms) > 1
}

/* Initiative Checks */

// InitiativeCommunityExists Are there any Initiative Communities?
func InitiativeCommunityExists(env environment, teamID models.TeamID, userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InitiativeCommunityExists")
	communities, err2 := strategy.StrategyCommunitiesDAOReadByPlatformID(teamID, 
		env.strategyCommunitiesTable)
	if err2 != nil {
		log.Printf("Failed InitiativeCommunityExists: %+v\n", err2)
	}
	for _, i := range communities {
		if i.Community == community.Initiative {
			return true
		}
	}
	return false
}

// InitiativesExistForMe Is the user an Advocate for any initiatives?
func InitiativesExistForMe(env environment, teamID models.TeamID, userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InitiativesExistForMe")
	inits := strategy.UserAdvocacyInitiatives(userID, 
		env.userObjectivesTable, env.userObjectivesTypeIndex, 0)
	if logEnabled {
		log.Println("Checked InitiativesExistForMe: ", len(inits))
	}
	return len(inits) > 0
}

// InitiativesExistInMyCapabilityCommunities Are there any Initiatives aligned with Capability Communities that I am in?
func InitiativesExistInMyCapabilityCommunities(env environment, teamID models.TeamID, userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InitiativesExistInMyCapabilityCommunities")
	var inits []models.StrategyInitiative
	if community.IsUserInCommunity(userID, env.communityUsersTable, env.communityUsersUserCommunityIndex, community.Strategy) {
		// User is in strategy community, return all the Initiatives
		inits = strategy.AllOpenStrategyInitiatives(teamID, env.initiativesTable, env.initiativesPlatformIndex,
			env.userObjectivesTable)
	} else {
		conn := env.connGen.ForPlatformID(teamID.ToPlatformID())
		inits = strategy.UserCapabilityCommunityInitiatives(userID, 
			env.strategyObjectivesTableName, env.strategyObjectivesPlatformIndex,
			env.initiativesTable, env.strategyInitiativesInitiativeCommunityIndex, 
			env.userObjectivesTable, env.communityUsersTable,
			env.communityUsersUserCommunityIndex, env.communityUsersUserIndex, conn)
	}
	if logEnabled {
		log.Println("Checked InitiativesExistInMyCapabilityCommunities: ", len(inits))
	}
	return len(inits) > 0
}

// InitiativesExistInMyInitiativeCommunities Are there any Initiatives aligned with Initiative Communities that I am in?
func InitiativesExistInMyInitiativeCommunities(env environment, teamID models.TeamID, userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InitiativesExistInMyInitiativeCommunities")
	inits := strategy.UserInitiativeCommunityInitiatives(userID,
		env.initiativesTable, env.strategyInitiativesInitiativeCommunityIndex,
		env.communityUsersTable, env.communityUsersUserIndex)
	if logEnabled {
		log.Println("Checked InitiativesExistInMyInitiativeCommunities: ", len(inits))
	}
	return len(inits) > 0
}

// StaleInitiativesExistForMe Is the user an Advocate for any initiatives
// that haven't been updated within the last month?
func StaleInitiativesExistForMe(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("StaleInitiativesExistForMe")
	initiativeObjs := strategy.UserInitiativesWithNoProgressInAWeek(userID, date,
		env.userObjectivesTable, env.userObjectivesUserIndex, env.userObjectivesProgressTable, 0)
	if logEnabled {
		log.Println("Checked StaleInitiativesExistForMe: ", len(initiativeObjs))
	}
	return len(initiativeObjs) > 0
}

// InitiativesDueInAWeek Are there any open Initiatives that exist for the user that are due in exactly 7 days
func InitiativesDueInAWeek(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("InitiativesDueInAWeek")
	op := strategy.InitiativesDueInAWeek(userID, date, 
		env.userObjectivesTable, env.userObjectivesUserIndex)
	if logEnabled {
		log.Println("Checked InitiativesDueInAWeek: ", len(op))
	}
	return len(op) > 0
}

// InitiativesDueInAMonth Are there any open Initiatives that exist for the user that are due in exactly 30 days
func InitiativesDueInAMonth(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("InitiativesDueInAMonth")
	op := strategy.InitiativesDueInAMonth(userID, date, 
		env.userObjectivesTable, env.userObjectivesUserIndex)
	if logEnabled {
		log.Println("Checked InitiativesDueInAMonth: ", len(op))
	}
	return len(op) > 0
}

// Are there any open Initiatives that exist for the user that are due in exactly 90 days
func InitiativesDueInAQuarter(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	op := strategy.InitiativesDueInAQuarter(userID, date, 
		env.userObjectivesTable, env.userObjectivesUserIndex)
	if logEnabled {
		log.Println("Checked InitiativesDueInAQuarter: ", len(op))
	}
	return len(op) > 0
}

/* Initiative Community Checks */

// InitiativeCommunityExistsForMe An Initiative Community exists for a
// objective community that the user is in.
func InitiativeCommunityExistsForMe(env environment, teamID models.TeamID,userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InitiativeCommunityExistsForMe")
	initComms := strategy.UserStrategyInitiativeCommunities(userID, 
		env.communityUsersTable, env.communityUsersUserCommunityIndex,
		env.communityUsersUserIndex, env.initiativeCommunitiesTableName, 
		env.initiativeCommunitiesPlatformIndex, env.strategyCommunitiesTable, teamID)
	if logEnabled {
		log.Println("Checked InitiativeCommunityExistsForMe: ", len(initComms))
	}
	return len(initComms) > 0
}

/* Miscellaneous Checks */

/* Team Values Check */

// TeamValuesExist Team values exist
func TeamValuesExist(env environment, teamID models.TeamID,userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("TeamValuesExist")
	vals := values.PlatformValues(teamID)
	return len(vals) > 0
}

// InCompetenciesCommunity The user is in the Values community
func InCompetenciesCommunity(env environment, teamID models.TeamID,userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InCompetenciesCommunity")
	return community.IsUserInCommunity(userID, 
		env.communityUsersTable, env.communityUsersUserCommunityIndex, community.Competency)
}

/* Holidays Check */

// HolidaysExist Holidays exist
func HolidaysExist(env environment, teamID models.TeamID,userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("HolidaysExist")
	platformID := teamID.ToPlatformID()
	conn := daosCommon.CreateConnectionFromEnv(platformID)
	vals := eholidays.AllUnsafe(conn)
	return len(vals) > 0
}

// InHRCommunity The user is in the HR user community
func InHRCommunity(env environment, teamID models.TeamID,userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InHRCommunity")
	return community.IsUserInCommunity(userID, 
		env.communityUsersTable, env.communityUsersUserCommunityIndex, community.HR)
}

/* Undelivered engagements check */

// UndeliveredEngagementsExistForMe There are undelivered engagements for the user
func UndeliveredEngagementsExistForMe(env environment, teamID models.TeamID,userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("UndeliveredEngagementsExistForMe")
	engs := user.NotPostedUnansweredNotIgnoredEngagements(userID, 
		env.engagementsTable, env.engagementsAnsweredIndex)
	return len(engs) > 0
}

// PostponedEventsExistForMe -
func PostponedEventsExistForMe(env environment, teamID models.TeamID,userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("PostponedEventsExistForMe")
	dao := postponedEvent.NewDAO(common.DeprecatedGetGlobalDns().Dynamo, "PostponedEventsExistForMe", env.clientID)

	events, err2 := dao.ReadByUserID(userID) //, engagementsTable, engagementsAnsweredIndex)
	if err2 != nil {
		log.Printf("PostponedEventsExistForMe user %s: %v\n", userID, err2)
	}
	res = err2 == nil && len(events) > 0

	return
}

// UndeliveredEngagementsOrPostponedEventsExistForMe -
func UndeliveredEngagementsOrPostponedEventsExistForMe(env environment, teamID models.TeamID,userID string, date business_time.Date) (res bool) {
	return UndeliveredEngagementsExistForMe(env, teamID, userID, date) ||
		PostponedEventsExistForMe(env, teamID, userID, date)
}

/* Reports exist check */

// ReportExists A performance report exists for the user
func ReportExists(env environment, teamID models.TeamID,userID string, dat business_time.Date) (res bool) {
	defer core.RecoverAsLogErrorf("ReportExists(userID=%s)", userID)
	key := coaching.UserReportIDForPreviousQuarter(dat.DateToTimeMidnight(), userID)
	res = common.DeprecatedGetGlobalS3().ObjectExists(env.reportsBucket, key)
	if logEnabled {
		log.Printf("Checked ReportExists(%s, %v): %v\n", userID, dat, res)
	}
	return
}

// FeedbackGivenForTheQuarter -
func FeedbackGivenForTheQuarter(env environment, teamID models.TeamID,userID string, date business_time.Date) (res bool) {
	defer core.RecoverAsLogErrorf("FeedbackGivenForTheQuarter(userID=%s)", userID)
	q := date.GetQuarter()
	y := date.GetYear()
	feedbacks, err2 := coaching.FeedbackGivenForTheQuarter(userID, q, y, 
		env.userFeedbackTable, 
		env.userFeedbackSourceQYIndex)
	if err2 != nil {
		log.Printf("Error with querying feedback given by the user %s: %v\n", userID, err2)
	}
	return len(feedbacks) > 0
}

// FeedbackForThePreviousQuarterExists -
func FeedbackForThePreviousQuarterExists(env environment, 
	teamID models.TeamID, userID string, 
	date business_time.Date,
) (res bool) {
	defer core.RecoverAsLogErrorf("FeedbackForThePreviousQuarterExists(userID=%s)", userID)
	conn := daosCommon.CreateConnectionGenFromEnv().ForPlatformID(teamID.ToPlatformID())
	q := date.GetPreviousQuarter()
	y := date.GetPreviousQuarterYear()
	var feedbacks []models.UserFeedback
	var err2 error
	feedbacks, err2 = coaching.FeedbackReceivedForTheQuarter(userID, q, y)(conn)
	res = len(feedbacks) > 0
	if err2 != nil {
		log.Printf("Error with querying feedback received by the user %s: %+v\n", userID, err2)
	}
	return
}

func CollectionNonEmpty(items []interface{}) (res bool) {
	return len(items) > 0
}

func CoacheesExistLogic(coacheeObjectives []userObjective.UserObjective) (res bool) {
	return len(coacheeObjectives) > 0
}

func AdvocatesExistLogic(advocateObjectives []userObjective.UserObjective) (res bool) {
	return len(advocateObjectives) > 0
}

func filterObjectivesByObjectiveType(objectives []userObjective.UserObjective, objectiveType userObjective.DevelopmentObjectiveType) (res []userObjective.UserObjective) {
	for _, objective := range objectives {
		if objective.ObjectiveType == objectiveType {
			res = append(res, objective)
		}
	}
	return
}

func LoadCoacheeObjectivesUnsafe(env environment, teamID models.TeamID, coachID string) (coacheeObjectives []userObjective.UserObjective) {
	conn := env.connGen.ForPlatformID(teamID.ToPlatformID())
	objectives := userObjective.ReadByAccountabilityPartnerUnsafe(coachID)(conn)
	coacheeObjectives = filterObjectivesByObjectiveType(objectives, userObjective.IndividualDevelopmentObjective)
	return
}

func LoadAdvocateeObjectivesUnsafe(coachID string) func (conn daosCommon.DynamoDBConnection) (advocateeObjectives []userObjective.UserObjective) {
	return func (conn daosCommon.DynamoDBConnection) (advocateeObjectives []userObjective.UserObjective) {
		objectives := userObjective.ReadByAccountabilityPartnerUnsafe(coachID)(conn)
		advocateeObjectives1 := filterObjectivesByObjectiveType(objectives, userObjective.StrategyDevelopmentObjective)
		advocateeObjectives2 := filterObjectivesByObjectiveType(objectives, userObjective.StrategyDevelopmentObjectiveIssue)
		advocateeObjectives = append(advocateeObjectives1, advocateeObjectives2...)
		return
	}
}

func CoacheesExist(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("CoacheesExist")
	coachID := userID
	conn := env.connGen.ForPlatformID(teamID.ToPlatformID())
	objectives := userObjective.ReadByAccountabilityPartnerUnsafe(coachID)(conn)
	coacheeObjectives := filterObjectivesByObjectiveType(objectives, userObjective.IndividualDevelopmentObjective)
	return CoacheesExistLogic(coacheeObjectives)
}

func AdvocatesExist(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("AdvocatesExist")
	conn := env.connGen.ForPlatformID(teamID.ToPlatformID())
	advocateObjectives := LoadAdvocateeObjectivesUnsafe(userID)(conn)
	return AdvocatesExistLogic(advocateObjectives)
}

func CanBeNudgedForIDOCreation(env environment, teamID models.TeamID, userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("CanBeNudgedForIDOCreation")
	inUserCommunity := community.IsUserInCommunity(userID, 
		env.communityUsersTable, env.communityUsersUserCommunityIndex, community.User)
	inInitiativeCommunity := InitiativeCommunityExistsForMe(env, teamID, userID, date)
	res = inUserCommunity || inInitiativeCommunity
	if logEnabled {
		log.Printf("User %s nudge for IDO creation: %v\n", userID, res)
	}
	return
}
