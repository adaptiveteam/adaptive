package adaptive_checks

import (
	"fmt"
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
	core_utils_go "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
	"github.com/adaptiveteam/adaptive/daos/postponedEvent"
)

/* IDO Checks */

// IDOsExistForMe Are there any IDO's that exist for the user?
func IDOsExistForMe(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("IDOsExistForMe")
	log.Println("Checking IDOsExistForMe")
	objs := objectives.AllUserObjectives(userID, userObjectivesTable, userObjectivesUserIndex,
		models.IndividualDevelopmentObjective, 0)
	res = len(objs) > 0
	log.Printf("IDOsExistForMe(%s, _): %v\n", userID, res)
	return
}

// IDOsDueInAWeek Are there any open IDO's that exist for the user that are due in exactly 7 days
func IDOsDueInAWeek(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("IDOsDueInAWeek")
	op := objectives.IDOsDueInAWeek(userID, date, userObjectivesTable, userObjectivesUserIndex)
	res = len(op) > 0
	log.Printf("IDOsDueInAWeek(%s, %v): %v\n", userID, date, res)
	return
}

// IDOsDueInAMonth Are there any open IDO's that exist for the user that are due in exactly in 30 days
func IDOsDueInAMonth(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("IDOsDueInAMonth")
	op := objectives.IDOsDueInAMonth(userID, date, userObjectivesTable, userObjectivesUserIndex)
	res = len(op) > 0
	log.Printf("IDOsDueInAMonth(%s, %v): %v\n", userID, date, res)
	return
}

// IDOsDueInAQuarter Are there any open IDO's that exist for the user that are due in exactly 90 days
func IDOsDueInAQuarter(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("IDOsDueInAQuarter")
	op := objectives.IDOsDueInAQuarter(userID, date, userObjectivesTable, userObjectivesUserIndex)
	res = len(op) > 0
	log.Printf("IDOsDueInAQuarter(%s, %v): %v\n", userID, date, res)
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
	log.Printf("StaleIDOsExist(%s, %v): %v\n", userID, date, res)
	return
}

/* Vision Checks */

// CompanyVisionExists Does the company vision exist?
func CompanyVisionExists(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("CompanyVisionExists")
	platformID := strategy.UserIDToPlatformID(userDAO)(userID)
	return strategy.StrategyVision(platformID, visionTable) != nil
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
	log.Println("Checked ObjectivesExistForMe: ", objs)
	return len(objs) > 0
}

// ObjectivesExist returns all the objectives for the organization
func ObjectivesExist(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("ObjectivesExist")
	log.Printf("Checking ObjectivesExist for userID=%s, date=%v\n", userID, date)
	objs := strategy.UserStrategyObjectives(userID, strategyObjectivesTableName, strategyObjectivesPlatformIndex,
		userObjectivesTable, communityUsersTable, communityUsersUserCommunityIndex)
	log.Println("Checked ObjectivesExist: ", len(objs))
	return len(objs) > 0
}

// Is the user the advocate for any objectives?
// that have no been updated within the last month?
func StaleObjectivesExistForMe(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("StaleObjectivesExistForMe")
	stratObjs := strategy.UserCapabilityObjectivesWithNoProgressInAMonth(userID, date,
		userObjectivesTable, userObjectivesUserIndex, userObjectivesProgressTable, 0)
	log.Println("Checked StaleObjectivesExistForMe: ", len(stratObjs))
	return len(stratObjs) > 0
}

// ObjectivesExistInMyCapabilityCommunities checks
//   if the user belong to any capability communities that have
// Capability Objectives allocated to them?
func ObjectivesExistInMyCapabilityCommunities(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("ObjectivesExistInMyCapabilityCommunities")
	log.Printf("Checking ObjectivesExistInMyCapabilityCommunities for userID=%s, date=%v\n", userID, date)
	objs := strategy.UserCommunityObjectives(userID,
		strategyObjectivesTableName, strategyObjectivesPlatformIndex,
		userObjectivesTable,
		communityUsersTable, communityUsersUserIndex)
	log.Printf("Checked ObjectivesExistInMyCapabilityCommunities: %d\n", len(objs))
	return len(objs) > 0
}

// CapabilityObjectivesDueInAWeek Are there any open Objectives that exist for the user that are due in exactly 7 days
func CapabilityObjectivesDueInAWeek(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("CapabilityObjectivesDueInAWeek")
	op := strategy.CapabilityObjectivesDueInAWeek(userID, date, userObjectivesTable, userObjectivesUserIndex)
	log.Println("Checked CapabilityObjectivesDueInAWeek: ", len(op))
	return len(op) > 0
}

// CapabilityObjectivesDueInAMonth Are there any open Objectives that exist for the user that are due in exactly 30 days
func CapabilityObjectivesDueInAMonth(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("CapabilityObjectivesDueInAMonth")
	op := strategy.CapabilityObjectivesDueInAMonth(userID, date, userObjectivesTable, userObjectivesUserIndex)
	log.Println("Checked CapabilityObjectivesDueInAMonth: ", len(op))
	return len(op) > 0
}

// CapabilityObjectivesDueInAQuarter Are there any open Objectives that exist for the user that are due in exactly 90 days
func CapabilityObjectivesDueInAQuarter(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("CapabilityObjectivesDueInAQuarter")
	op := strategy.CapabilityObjectivesDueInAQuarter(userID, date, userObjectivesTable, userObjectivesUserIndex)
	log.Println("Checked CapabilityObjectivesDueInAQuarter: ", len(op))
	return len(op) > 0
}

/* Capabilitity Community Checks */

// InCapabilityCommunity Is the user in any Capability Community?
func InCapabilityCommunity(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InCapabilityCommunity")
	capComms, _ := strategy.UserCapabilityInitiativeCommunities(userID, communityUsersTable, communityUsersUserIndex)
	log.Println("Checked InCapabilityCommunity: ", len(capComms))
	return len(capComms) > 0
}

// CapabilityCommunityExists Does there exist a capabilility community?
func CapabilityCommunityExists(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("CapabilityCommunityExists")
	platformID := strategy.UserIDToPlatformID(userDAO)(userID)
	capComms := strategy.AllCapabilityCommunities(platformID, capabilityCommunitiesTable,
		capabilityCommunitiesPlatformIndex, strategyCommunitiesTable)
	log.Println("Checked CapabilityCommunityExists: ", len(capComms))
	return len(capComms) > 0
}

// MultipleCapabilityCommunitiesExists Is there more than one capability community?
func MultipleCapabilityCommunitiesExists(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("MultipleCapabilityCommunitiesExists")
	platformID := strategy.UserIDToPlatformID(userDAO)(userID)
	capComms := strategy.AllCapabilityCommunities(platformID, capabilityCommunitiesTable,
		capabilityCommunitiesPlatformIndex, strategyCommunitiesTable)
	log.Println("Checked MultipleCapabilityCommunitiesExists: ", len(capComms))
	return len(capComms) > 1
}

/* Initiative Checks */

// InitiativeCommunityExists Are there any Initiative Communities?
func InitiativeCommunityExists(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InitiativeCommunityExists")
	platformID := strategy.UserIDToPlatformID(userDAO)(userID)

	communities, err2 := strategy.StrategyCommunitiesDAOReadByPlatformID(platformID, strategyCommunitiesTable)
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
	log.Println("Checked InitiativesExistForMe: ", len(inits))
	return len(inits) > 0
}

// InitiativesExistInMyCapabilityCommunities Are there any Initiatives aligned with Capability Communities that I am in?
func InitiativesExistInMyCapabilityCommunities(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InitiativesExistInMyCapabilityCommunities")
	var inits []models.StrategyInitiative
	if community.IsUserInCommunity(userID, communityUsersTable, communityUsersUserCommunityIndex, community.Strategy) {
		// User is in strategy community, return all the Initiatives
		platformID := strategy.UserIDToPlatformID(userDAO)(userID)
		inits = strategy.AllOpenStrategyInitiatives(platformID, initiativesTable, initiativesPlatformIndex,
			userObjectivesTable)
	} else {
		inits = strategy.UserCapabilityCommunityInitiatives(userID, strategyObjectivesTableName, strategyObjectivesPlatformIndex,
			initiativesTable, strategyInitiativesInitiativeCommunityIndex, userObjectivesTable, communityUsersTable,
			communityUsersUserCommunityIndex, communityUsersUserIndex)
	}
	log.Println("Checked InitiativesExistInMyCapabilityCommunities: ", len(inits))
	return len(inits) > 0
}

// InitiativesExistInMyInitiativeCommunities Are there any Initiatives aligned with Initiative Communities that I am in?
func InitiativesExistInMyInitiativeCommunities(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InitiativesExistInMyInitiativeCommunities")
	inits := strategy.UserInitiativeCommunityInitiatives(userID,
		initiativesTable, strategyInitiativesInitiativeCommunityIndex,
		communityUsersTable, communityUsersUserIndex)
	log.Println("Checked InitiativesExistInMyInitiativeCommunities: ", len(inits))
	return len(inits) > 0
}

// StaleInitiativesExistForMe Is the user an Advocate for any initiatives
// that haven't been updated within the last month?
func StaleInitiativesExistForMe(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("StaleInitiativesExistForMe")
	initiativeObjs := strategy.UserInitiativesWithNoProgressInAWeek(userID, date,
		userObjectivesTable, userObjectivesUserIndex, userObjectivesProgressTable, 0)
	log.Println("Checked StaleInitiativesExistForMe: ", len(initiativeObjs))
	return len(initiativeObjs) > 0
}

// InitiativesDueInAWeek Are there any open Initiatives that exist for the user that are due in exactly 7 days
func InitiativesDueInAWeek(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("InitiativesDueInAWeek")
	op := strategy.InitiativesDueInAWeek(userID, date, userObjectivesTable, userObjectivesUserIndex)
	log.Println("Checked InitiativesDueInAWeek: ", len(op))
	return len(op) > 0
}

// InitiativesDueInAMonth Are there any open Initiatives that exist for the user that are due in exactly 30 days
func InitiativesDueInAMonth(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("InitiativesDueInAMonth")
	op := strategy.InitiativesDueInAMonth(userID, date, userObjectivesTable, userObjectivesUserIndex)
	log.Println("Checked InitiativesDueInAMonth: ", len(op))
	return len(op) > 0
}

// Are there any open Initiatives that exist for the user that are due in exactly 90 days
func InitiativesDueInAQuarter(userID string, date business_time.Date) (res bool) {
	op := strategy.InitiativesDueInAQuarter(userID, date, userObjectivesTable, userObjectivesUserIndex)
	log.Println("Checked InitiativesDueInAQuarter: ", len(op))
	return len(op) > 0
}

/* Initiative Community Checks */

// InitiativeCommunityExistsForMe An Initiative Community exists for a
// capability community that the user is in.
func InitiativeCommunityExistsForMe(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("InitiativeCommunityExistsForMe")
	platformID := strategy.UserIDToPlatformID(userDAO)(userID)
	initComms := strategy.UserStrategyInitiativeCommunities(userID, communityUsersTable, communityUsersUserCommunityIndex,
		communityUsersUserIndex, initiativeCommunitiesTableName, initiativeCommunitiesPlatformIndex, strategyCommunitiesTable, platformID)
	log.Println("Checked InitiativeCommunityExistsForMe: ", len(initComms))
	return len(initComms) > 0
}

/* Miscellaneous Checks */

/* Team Values Check */

// TeamValuesExist Team values exist
func TeamValuesExist(userID string, _ business_time.Date) (res bool) {
	defer RecoverToLog("TeamValuesExist")
	platformID := UserIDToPlatformID(userDAO)(userID)
	vals := values.PlatformValues(platformID)
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
	platformID := UserIDToPlatformID(userDAO)(userID)
	vals := adHocHolidaysTableDao.ForPlatformID(platformID).AllUnsafe()
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
	dao := postponedEvent.NewDAO(common.DeprecatedGetGlobalDns().Dynamo,"PostponedEventsExistForMe", clientID)
	
	events, err2 := dao.ReadByUserID(userID)//, engagementsTable, engagementsAnsweredIndex)
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
	defer RecoverToLog("ReportExists")
	key, err := coaching.UserReportIDForPreviousQuarter(models.UserEngage{
		UserId:   userID,
		Date:     dat.DateToString(string(core_utils_go.ISODateLayout)),
		OnDemand: false,
	})
	if err == nil {
		res = common.DeprecatedGetGlobalS3().ObjectExists(reportsBucket, key)
	}
	log.Printf("Checked ReportExists(%s, %v): %v\n", userID, dat, res)
	return
}

// FeedbackGivenForTheQuarter -
func FeedbackGivenForTheQuarter(userID string, date business_time.Date) (res bool) {
	defer RecoverToLog("FeedbackGivenForTheQuarter")
	q := date.GetQuarter()
	y := date.GetYear()
	feedbacks, err := coaching.FeedbackGivenForTheQuarter(userID, q, y, userFeedbackTable, userFeedbackSourceQYIndex)
	if err != nil {
		log.Println(fmt.Sprintf("Error with querying feedback given by the user %s", userID))
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
	log.Println(fmt.Sprintf("User %s nudge for IDO creation: %v", userID, res))
	return
}
