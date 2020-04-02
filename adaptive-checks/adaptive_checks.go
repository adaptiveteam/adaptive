package adaptive_checks

import (
	"log"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/values"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	core_utils_go "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/postponedEvent"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
)

/* IDO Checks */

const logEnabled = false

// IDOsExistForMe Are there any IDO's that exist for the user?
func IDOsExistForMe(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("IDOsExistForMe")
	if logEnabled {
		log.Println("Checking IDOsExistForMe")
	}
	objs := objectives.AllUserObjectives(userID, userObjectivesTable, userObjectivesUserIndex,
		models.IndividualDevelopmentObjective, 0)
	res = len(objs) > 0
	if logEnabled {
		log.Printf("IDOsExistForMe(%s, _): %v\n", userID, res)
	}
	return
}

// IDOsDueInAWeek Are there any open IDO's that exist for the user that are due in exactly 7 days
func IDOsDueInAWeek(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("IDOsDueInAWeek")
	op := objectives.IDOsDueInAWeek(userID, date, userObjectivesTable, userObjectivesUserIndex)
	res = len(op) > 0
	if logEnabled {
		log.Printf("IDOsDueInAWeek(%s, %v): %v\n", userID, date, res)
	}
	return
}

// IDOsDueInAMonth Are there any open IDO's that exist for the user that are due in exactly in 30 days
func IDOsDueInAMonth(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("IDOsDueInAMonth")
	op := objectives.IDOsDueInAMonth(userID, date, userObjectivesTable, userObjectivesUserIndex)
	res = len(op) > 0
	if logEnabled {
		log.Printf("IDOsDueInAMonth(%s, %v): %v\n", userID, date, res)
	}
	return
}

// IDOsDueInAQuarter Are there any open IDO's that exist for the user that are due in exactly 90 days
func IDOsDueInAQuarter(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("IDOsDueInAQuarter")
	op := objectives.IDOsDueInAQuarter(userID, date, userObjectivesTable, userObjectivesUserIndex)
	res = len(op) > 0
	if logEnabled {
		log.Printf("IDOsDueInAQuarter(%s, %v): %v\n", userID, date, res)
	}
	return
}

// StaleIDOs checks if there are any stale IDO's for a user
func StaleIDOs(userID string, date business_time.Date) []models.UserObjective {
	defer RecoverToLog("StaleIDOs")
	return objectives.UserIDOsWithNoProgressInAWeek(userID, date, userObjectivesTable, userObjectivesUserIndex,
		userObjectivesProgressTable)
}

// StaleIDOsExist checks that there are IDOs without recent progress
func StaleIDOsExist(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("StaleIDOsExist")
	staleIDOs := StaleIDOs(userID, date)
	res = len(staleIDOs) > 0
	if logEnabled {
		log.Printf("StaleIDOsExist(%s, %v): %v\n", userID, date, res)
	}
	return
}

/* Vision Checks */

// CompanyVisionExists Does the company vision exist?
func CompanyVisionExists(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("CompanyVisionExists")
	teamID := strategy.UserIDToTeamID(userDAO)(userID)
	return strategy.StrategyVision(teamID, visionTable) != nil
}

// InStrategyCommunity Is the user in the strategy community
func InStrategyCommunity(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InStrategyCommunity")
	return community.IsUserInCommunity(userID, communityUsersTable, communityUsersUserCommunityIndex, community.Strategy)
}

/* Objective Checks */

// ObjectivesExistForMe Is the user the advocate for any objectives?
func ObjectivesExistForMe(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("ObjectivesExistForMe")
	objs := strategy.UserAdvocacyObjectives(userID, userObjectivesTable, userObjectivesTypeIndex, 0)
	if logEnabled {
		log.Println("Checked ObjectivesExistForMe: ", objs)
	}
	return len(objs) > 0
}

// ObjectivesExist returns all the objectives for the organization
func ObjectivesExist(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("ObjectivesExist")
	if logEnabled {
		log.Printf("Checking ObjectivesExist for userID=%s, date=%v\n", userID, date)
	}
	platformID := UserIDToPlatformID(userDAO)(userID)
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
func StaleObjectivesExistForMe(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("StaleObjectivesExistForMe")
	stratObjs := strategy.UserCapabilityObjectivesWithNoProgressInAMonth(userID, date,
		userObjectivesTable, userObjectivesUserIndex, userObjectivesProgressTable, 0)
	if logEnabled {
		log.Println("Checked StaleObjectivesExistForMe: ", len(stratObjs))
	}
	return len(stratObjs) > 0
}

// ObjectivesExistInMyCapabilityCommunities checks
//   if the user belong to any capability communities that have
// Capability Objectives allocated to them?
func ObjectivesExistInMyCapabilityCommunities(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("ObjectivesExistInMyCapabilityCommunities")
	if logEnabled {
		log.Printf("Checking ObjectivesExistInMyCapabilityCommunities for userID=%s, date=%v\n", userID, date)
	}
	objs := strategy.UserCommunityObjectives(userID,
		strategyObjectivesTableName, strategyObjectivesPlatformIndex,
		userObjectivesTable,
		communityUsersTable, communityUsersUserIndex)
	if logEnabled {
		log.Printf("Checked ObjectivesExistInMyCapabilityCommunities: %d\n", len(objs))
	}
	return len(objs) > 0
}

// CapabilityObjectivesDueInAWeek Are there any open Objectives that exist for the user that are due in exactly 7 days
func CapabilityObjectivesDueInAWeek(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("CapabilityObjectivesDueInAWeek")
	op := strategy.CapabilityObjectivesDueInAWeek(userID, date, userObjectivesTable, userObjectivesUserIndex)
	if logEnabled {
		log.Println("Checked CapabilityObjectivesDueInAWeek: ", len(op))
	}
	return len(op) > 0
}

// CapabilityObjectivesDueInAMonth Are there any open Objectives that exist for the user that are due in exactly 30 days
func CapabilityObjectivesDueInAMonth(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("CapabilityObjectivesDueInAMonth")
	op := strategy.CapabilityObjectivesDueInAMonth(userID, date, userObjectivesTable, userObjectivesUserIndex)
	if logEnabled {
		log.Println("Checked CapabilityObjectivesDueInAMonth: ", len(op))
	}
	return len(op) > 0
}

// CapabilityObjectivesDueInAQuarter Are there any open Objectives that exist for the user that are due in exactly 90 days
func CapabilityObjectivesDueInAQuarter(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("CapabilityObjectivesDueInAQuarter")
	op := strategy.CapabilityObjectivesDueInAQuarter(userID, date, userObjectivesTable, userObjectivesUserIndex)
	if logEnabled {
		log.Println("Checked CapabilityObjectivesDueInAQuarter: ", len(op))
	}
	return len(op) > 0
}

/* Capabilitity Community Checks */

// InCapabilityCommunity Is the user in any Objective Community?
func InCapabilityCommunity(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InCapabilityCommunity")
	capComms, _ := strategy.UserCapabilityInitiativeCommunities(userID, communityUsersTable, communityUsersUserIndex)
	if logEnabled {
		log.Println("Checked InCapabilityCommunity: ", len(capComms))
	}
	return len(capComms) > 0
}

// CapabilityCommunityExists Does there exist a capabilility community?
func CapabilityCommunityExists(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("CapabilityCommunityExists")
	teamID := strategy.UserIDToTeamID(userDAO)(userID)
	capComms := strategy.AllCapabilityCommunities(teamID, capabilityCommunitiesTable,
		capabilityCommunitiesPlatformIndex, strategyCommunitiesTable)
	if logEnabled {
		log.Println("Checked CapabilityCommunityExists: ", len(capComms))
	}
	return len(capComms) > 0
}

// MultipleCapabilityCommunitiesExists Is there more than one objective community?
func MultipleCapabilityCommunitiesExists(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("MultipleCapabilityCommunitiesExists")
	teamID := strategy.UserIDToTeamID(userDAO)(userID)
	capComms := strategy.AllCapabilityCommunities(teamID, capabilityCommunitiesTable,
		capabilityCommunitiesPlatformIndex, strategyCommunitiesTable)
	if logEnabled {
		log.Println("Checked MultipleCapabilityCommunitiesExists: ", len(capComms))
	}
	return len(capComms) > 1
}

/* Initiative Checks */

// InitiativeCommunityExists Are there any Initiative Communities?
func InitiativeCommunityExists(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InitiativeCommunityExists")
	teamID := strategy.UserIDToTeamID(userDAO)(userID)

	communities, err2 := strategy.StrategyCommunitiesDAOReadByPlatformID(teamID, strategyCommunitiesTable)
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
func InitiativesExistForMe(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InitiativesExistForMe")
	inits := strategy.UserAdvocacyInitiatives(userID, userObjectivesTable, userObjectivesTypeIndex, 0)
	if logEnabled {
		log.Println("Checked InitiativesExistForMe: ", len(inits))
	}
	return len(inits) > 0
}

// InitiativesExistInMyCapabilityCommunities Are there any Initiatives aligned with Capability Communities that I am in?
func InitiativesExistInMyCapabilityCommunities(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InitiativesExistInMyCapabilityCommunities")
	var inits []models.StrategyInitiative
	if community.IsUserInCommunity(userID, communityUsersTable, communityUsersUserCommunityIndex, community.Strategy) {
		// User is in strategy community, return all the Initiatives
		teamID := strategy.UserIDToTeamID(userDAO)(userID)
		inits = strategy.AllOpenStrategyInitiatives(teamID, initiativesTable, initiativesPlatformIndex,
			userObjectivesTable)
	} else {
		inits = strategy.UserCapabilityCommunityInitiatives(userID, strategyObjectivesTableName, strategyObjectivesPlatformIndex,
			initiativesTable, strategyInitiativesInitiativeCommunityIndex, userObjectivesTable, communityUsersTable,
			communityUsersUserCommunityIndex, communityUsersUserIndex)
	}
	if logEnabled {
		log.Println("Checked InitiativesExistInMyCapabilityCommunities: ", len(inits))
	}
	return len(inits) > 0
}

// InitiativesExistInMyInitiativeCommunities Are there any Initiatives aligned with Initiative Communities that I am in?
func InitiativesExistInMyInitiativeCommunities(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InitiativesExistInMyInitiativeCommunities")
	inits := strategy.UserInitiativeCommunityInitiatives(userID,
		initiativesTable, strategyInitiativesInitiativeCommunityIndex,
		communityUsersTable, communityUsersUserIndex)
	if logEnabled {
		log.Println("Checked InitiativesExistInMyInitiativeCommunities: ", len(inits))
	}
	return len(inits) > 0
}

// StaleInitiativesExistForMe Is the user an Advocate for any initiatives
// that haven't been updated within the last month?
func StaleInitiativesExistForMe(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("StaleInitiativesExistForMe")
	initiativeObjs := strategy.UserInitiativesWithNoProgressInAWeek(userID, date,
		userObjectivesTable, userObjectivesUserIndex, userObjectivesProgressTable, 0)
	if logEnabled {
		log.Println("Checked StaleInitiativesExistForMe: ", len(initiativeObjs))
	}
	return len(initiativeObjs) > 0
}

// InitiativesDueInAWeek Are there any open Initiatives that exist for the user that are due in exactly 7 days
func InitiativesDueInAWeek(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("InitiativesDueInAWeek")
	op := strategy.InitiativesDueInAWeek(userID, date, userObjectivesTable, userObjectivesUserIndex)
	if logEnabled {
		log.Println("Checked InitiativesDueInAWeek: ", len(op))
	}
	return len(op) > 0
}

// InitiativesDueInAMonth Are there any open Initiatives that exist for the user that are due in exactly 30 days
func InitiativesDueInAMonth(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("InitiativesDueInAMonth")
	op := strategy.InitiativesDueInAMonth(userID, date, userObjectivesTable, userObjectivesUserIndex)
	if logEnabled {
		log.Println("Checked InitiativesDueInAMonth: ", len(op))
	}
	return len(op) > 0
}

// Are there any open Initiatives that exist for the user that are due in exactly 90 days
func InitiativesDueInAQuarter(userID string, date business_time.Date) (res bool) {
	op := strategy.InitiativesDueInAQuarter(userID, date, userObjectivesTable, userObjectivesUserIndex)
	if logEnabled {
		log.Println("Checked InitiativesDueInAQuarter: ", len(op))
	}
	return len(op) > 0
}

/* Initiative Community Checks */

// InitiativeCommunityExistsForMe An Initiative Community exists for a
// objective community that the user is in.
func InitiativeCommunityExistsForMe(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InitiativeCommunityExistsForMe")
	teamID := strategy.UserIDToTeamID(userDAO)(userID)
	initComms := strategy.UserStrategyInitiativeCommunities(userID, communityUsersTable, communityUsersUserCommunityIndex,
		communityUsersUserIndex, initiativeCommunitiesTableName, initiativeCommunitiesPlatformIndex, strategyCommunitiesTable, teamID)
	if logEnabled {
		log.Println("Checked InitiativeCommunityExistsForMe: ", len(initComms))
	}
	return len(initComms) > 0
}

/* Miscellaneous Checks */

/* Team Values Check */

// TeamValuesExist Team values exist
func TeamValuesExist(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("TeamValuesExist")
	teamID := UserIDToTeamID(userDAO)(userID)
	vals := values.PlatformValues(teamID)
	return len(vals) > 0
}

// InCompetenciesCommunity The user is in the Values community
func InCompetenciesCommunity(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InCompetenciesCommunity")
	return community.IsUserInCommunity(userID, communityUsersTable, communityUsersUserCommunityIndex, community.Competency)
}

/* Holidays Check */

// HolidaysExist Holidays exist
func HolidaysExist(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("HolidaysExist")
	teamID := UserIDToTeamID(userDAO)(userID)
	vals := adHocHolidaysTableDao.ForPlatformID(teamID).AllUnsafe()
	return len(vals) > 0
}

// InHRCommunity The user is in the HR user community
func InHRCommunity(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InHRCommunity")
	return community.IsUserInCommunity(userID, communityUsersTable, communityUsersUserCommunityIndex, community.HR)
}

/* Undelivered engagements check */

// UndeliveredEngagementsExistForMe There are undelivered engagements for the user
func UndeliveredEngagementsExistForMe(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("UndeliveredEngagementsExistForMe")
	engs := user.NotPostedUnansweredNotIgnoredEngagements(userID, engagementsTable, engagementsAnsweredIndex)
	return len(engs) > 0
}

// PostponedEventsExistForMe -
func PostponedEventsExistForMe(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("PostponedEventsExistForMe")
	dao := postponedEvent.NewDAO(common.DeprecatedGetGlobalDns().Dynamo, "PostponedEventsExistForMe", clientID)

	events, err2 := dao.ReadByUserID(userID) //, engagementsTable, engagementsAnsweredIndex)
	if err2 != nil {
		log.Printf("PostponedEventsExistForMe user %s: %v\n", userID, err2)
	}
	res = err2 == nil && len(events) > 0

	return
}

// UndeliveredEngagementsOrPostponedEventsExistForMe -
func UndeliveredEngagementsOrPostponedEventsExistForMe(userID string, date business_time.Date) (res bool) {
	return UndeliveredEngagementsExistForMe(userID, date) ||
		PostponedEventsExistForMe(userID, date)
}

/* Reports exist check */

// ReportExists A performance report exists for the user
func ReportExists(userID string, dat business_time.Date) (res bool) {
	defer core.RecoverAsLogErrorf("ReportExists(userID=%s)", userID)
	key, err2 := coaching.UserReportIDForPreviousQuarter(models.UserEngage{
		UserID:   userID,
		Date:     dat.DateToString(string(core_utils_go.ISODateLayout)),
		OnDemand: false,
	})
	if err2 == nil {
		res = common.DeprecatedGetGlobalS3().ObjectExists(reportsBucket, key)
	} else {
		log.Printf("ReportExists user %s: %v\n", userID, err2)
	}
	if logEnabled {
		log.Printf("Checked ReportExists(%s, %v): %v\n", userID, dat, res)
	}
	return
}

// FeedbackGivenForTheQuarter -
func FeedbackGivenForTheQuarter(userID string, date business_time.Date) (res bool) {
	defer core.RecoverAsLogErrorf("FeedbackGivenForTheQuarter(userID=%s)", userID)
	q := date.GetQuarter()
	y := date.GetYear()
	feedbacks, err2 := coaching.FeedbackGivenForTheQuarter(userID, q, y, userFeedbackTable, userFeedbackSourceQYIndex)
	if err2 != nil {
		log.Printf("Error with querying feedback given by the user %s: %v\n", userID, err2)
	}
	return len(feedbacks) > 0
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

func LoadCoacheeObjectivesUnsafe(coachID string) (coacheeObjectives []userObjective.UserObjective) {
	objectives := userObjectiveDAO.ReadByAccountabilityPartnerUnsafe(coachID)
	coacheeObjectives = filterObjectivesByObjectiveType(objectives, userObjective.IndividualDevelopmentObjective)
	return
}

func LoadAdvocateeObjectivesUnsafe(coachID string) (advocateeObjectives []userObjective.UserObjective) {
	objectives := userObjectiveDAO.ReadByAccountabilityPartnerUnsafe(coachID)
	advocateeObjectives = filterObjectivesByObjectiveType(objectives, userObjective.StrategyDevelopmentObjective)
	return
}

func CoacheesExist(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("CoacheesExist")
	coachID := userID
	objectives := userObjectiveDAO.ReadByAccountabilityPartnerUnsafe(coachID)
	coacheeObjectives := filterObjectivesByObjectiveType(objectives, userObjective.IndividualDevelopmentObjective)
	return CoacheesExistLogic(coacheeObjectives)
}

func AdvocatesExist(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("AdvocatesExist")
	objectives := userObjectiveDAO.ReadByAccountabilityPartnerUnsafe(userID)
	advocateObjectives := filterObjectivesByObjectiveType(objectives, userObjective.StrategyDevelopmentObjective)
	return AdvocatesExistLogic(advocateObjectives)
}

func CanBeNudgedForIDOCreation(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("CanBeNudgedForIDOCreation")
	inUserCommunity := community.IsUserInCommunity(userID, communityUsersTable, communityUsersUserCommunityIndex, community.User)
	inInitiativeCommunity := InitiativeCommunityExistsForMe(userID, date)
	res = inUserCommunity || inInitiativeCommunity
	if logEnabled {
		log.Printf("User %s nudge for IDO creation: %v\n", userID, res)
	}
	return
}
