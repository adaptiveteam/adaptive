package strategy

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"strings"
)

// Constructing Dynamo query expression based on index for the platform
func platformIndexExpr(index string, platformID models.PlatformID) awsutils.DynamoIndexExpression {
	return awsutils.DynamoIndexExpression{
		IndexName: index,
		Condition: "platform_id = :p",
		Attributes: map[string]interface{}{
			":p": string(platformID),
		},
	}
}

// UserStrategyObjectives returns all open objectives associated with a user
// If user is in strategy community, we return all objectives
// Else we return those objectives associated with capability communities that the user is a part of
// USED
func UserStrategyObjectives(userID string,
	strategyObjectivesTable, strategyObjectivesPlatformIndex, userObjectivesTable string,
	communityUsersTable, communityUsersUserCommunityIndex string) []models.StrategyObjective {
	log.Printf("UserStrategyObjectives(userID=%s, strategyObjectivesTable=%s, strategyObjectivesPlatformIndex=%s, userObjectivesTable=%s, communityUsersTable=%s, communityUsersUserCommunityIndex=%s)", 
									   userID,    strategyObjectivesTable,    strategyObjectivesPlatformIndex,    userObjectivesTable,    communityUsersTable,    communityUsersUserCommunityIndex)
	if community.IsUserInCommunity(userID, communityUsersTable, communityUsersUserCommunityIndex, community.Strategy) {
		log.Println(fmt.Sprintf("### User %s is in strategy community, showing all the objectives", userID))
		platformID := UserIDToPlatformID(userDAO())(userID)
		return AllOpenStrategyObjectives(platformID, strategyObjectivesTable, strategyObjectivesPlatformIndex,
			userObjectivesTable)
	} else {
		log.Println(fmt.Sprintf("### User %s is not in strategy community, showing relevant objectives", userID))
		return UserCommunityObjectives(userID, strategyObjectivesTable, strategyObjectivesPlatformIndex, userObjectivesTable,
			communityUsersTable, communityUsersUserCommunityIndex)
	}
}

// USED
func allStrategyObjectives(platformID models.PlatformID, strategyObjectivesTable,
	strategyObjectivesPlatformIndex string) []models.StrategyObjective {
	var objs []models.StrategyObjective
	err := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(strategyObjectivesTable,
		platformIndexExpr(strategyObjectivesPlatformIndex, platformID),
		map[string]string{}, true, -1, &objs)
	core.ErrorHandler(err, common.DeprecatedGetGlobalDns().Namespace, fmt.Sprintf("Could not query %s index on %s table",
		strategyObjectivesPlatformIndex,
		strategyObjectivesTable))
	return objs
}

// AllOpenStrategyObjectives returns a slice of open strategy objectives: capability, customer and financial objectives
// USED
func AllOpenStrategyObjectives(platformID models.PlatformID, strategyObjectivesTable, strategyObjectivesPlatformIndex, 
	userObjectivesTable string) []models.StrategyObjective {
	allObjs := allStrategyObjectives(platformID, strategyObjectivesTable, strategyObjectivesPlatformIndex)
	log.Printf("AllOpenStrategyObjectives: len(allObjs)=%d\n", len(allObjs))

	var res []models.StrategyObjective
	for _, each := range allObjs {
		// there has to be at least one capability community id
		// TODO: This presents a tricky scenario when original capability community is updated. Think about this.
		// Customer and financial objectives have no capability communities associated with them. For them,we only use the ID
		id := each.ID
		// if len(each.CapabilityCommunityIDs) > 0 {
		// 	id = fmt.Sprintf("%s_%s", each.ID, each.CapabilityCommunityIDs[0])
		// }
		userObj, err2 := getUserObjectiveByID(userObjectivesTable, id)
		if err2 == nil {
			log.Printf("AllOpenStrategyObjectives: userObj(id=%s).Completed=%d\n", id, userObj.Completed)
			if userObj.Completed == 0 {
				res = append(res, each)
			}
		} else {
			log.Printf("AllOpenStrategyObjectives, error for userObj(id=%s) %v", id, err2)
		}
	}
	return res
}

func getUserObjectiveByID(userObjectivesTable string, id string) (uo models.UserObjective, err error) {
	defer recoverToErrorVar("getUserObjectiveByID", &err)
	res := objectives.UserObjectiveById(userObjectivesTable, id, common.DeprecatedGetGlobalDns())
	if res == nil {
		err = fmt.Errorf("Not found userObj(id=%s)", id)
	} else {
		uo = *res
	}
	return
}
func AllStrategyInitiatives(platformID models.PlatformID, strategyInitiativesTable,
	strategyInitiativesPlatformIndex string) []models.StrategyInitiative {
	var sis []models.StrategyInitiative
	err := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(strategyInitiativesTable,
		platformIndexExpr(strategyInitiativesPlatformIndex, platformID),
		map[string]string{}, true, -1, &sis)
	core.ErrorHandler(err, common.DeprecatedGetGlobalDns().Namespace, fmt.Sprintf("Could not query %s index on %s table",
		strategyInitiativesPlatformIndex,
		strategyInitiativesTable))
	return sis
}

// AllOpenStrategyInitiatives returns a slice of all open strategy initiatives
func AllOpenStrategyInitiatives(platformID models.PlatformID, strategyInitiativesTable, strategyInitiativesPlatformIndex, 
	userObjectivesTable string) []models.StrategyInitiative {
	inits := AllStrategyInitiatives(platformID, strategyInitiativesTable, strategyInitiativesPlatformIndex)
	var res []models.StrategyInitiative
	for _, each := range inits {
		userObj := objectives.UserObjectiveById(userObjectivesTable, each.ID, common.DeprecatedGetGlobalDns())
		if userObj != nil && userObj.Completed == 0 {
			res = append(res, each)
		}
	}
	return res
}

// StrategyCommunityOrEmptyByID retrives the strategy community based on the id of the community
func StrategyCommunityOrEmptyByID(id, strategyCommunitiesTable string) (comm StrategyCommunity, found bool, err error) {
	params := map[string]*dynamodb.AttributeValue{
		"id": dynString(id),
	}
	found, err = d().GetItemOrEmptyFromTable(strategyCommunitiesTable, params, &comm)
	return
}

// StrategyCommunityByID retrives the strategy community based on the id of the community
func StrategyCommunityByID(id, strategyCommunitiesTable string) StrategyCommunity {
	return StrategyCommunityByIDUnsafe(id, strategyCommunitiesTable)
}

// StrategyCommunityByIDUnsafe retrives the strategy community based on the id of the community
func StrategyCommunityByIDUnsafe(id, strategyCommunitiesTable string) StrategyCommunity {
	params := map[string]*dynamodb.AttributeValue{
		"id": dynString(id),
	}
	var comm StrategyCommunity
	err2 := d().GetItemFromTable(strategyCommunitiesTable, params, &comm)
	core.ErrorHandler(err2, namespace(), fmt.Sprintf("StrategyCommunityByIDUnsafe: Could not find %s in %s table", id, strategyCommunitiesTable))
	return comm
}

// AllCapabilityCommunities Get all the capability communities, 
// that have Adaptive associated, for the platform ID
func AllCapabilityCommunities(platformID models.PlatformID,
	capabilityCommunitiesTable, capabilityCommunitiesPlatformIndex, strategyCommunitiesTable string) []CapabilityCommunity {
	var ccs, op []CapabilityCommunity
	err2 := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(capabilityCommunitiesTable,
		platformIndexExpr(capabilityCommunitiesPlatformIndex, platformID),
		map[string]string{}, true, -1, &ccs)
	core.ErrorHandler(err2, "AllCapabilityCommunities.all capabilityCommunities for platformID", fmt.Sprintf("Could not query %s index on %s table. platformID=%s",
		capabilityCommunitiesPlatformIndex,
		capabilityCommunitiesTable,
		platformID))

	for _, each := range ccs {
		stratComm, found, err3 := StrategyCommunityOrEmptyByID(each.ID, strategyCommunitiesTable)
		core.ErrorHandler(err3, "AllCapabilityCommunities.StrategyCommunityOrEmptyByID", "Could not StrategyCommunityOrEmptyByID")
		if found {
			if stratComm.ChannelCreated == 1 {
				op = append(op, each)
			}
		} else {
			log.Printf("Couldn't find the respective StrategyCommunity for id=%s", each.ID)
		}
	}
	return op
}

// AllStrategyInitiativeCommunities - Get all the initiative communities 
// for the platform ID
func AllStrategyInitiativeCommunities(platformID models.PlatformID, initiativeCommunitiesTable,
	initiativeCommunitiesPlatformIndex, strategyCommunitiesTable string) []StrategyInitiativeCommunity {
	var sics, op []StrategyInitiativeCommunity
	err2 := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(
		initiativeCommunitiesTable,
		platformIndexExpr(initiativeCommunitiesPlatformIndex, platformID),
		map[string]string{}, true, -1, &sics)
	core.ErrorHandler(err2, "AllStrategyInitiativeCommunities", fmt.Sprintf("Could not query %s index on %s table",
		initiativeCommunitiesPlatformIndex, initiativeCommunitiesTable))
	for _, each := range sics {
		stratComm, found, err3 := StrategyCommunityOrEmptyByID(each.ID, strategyCommunitiesTable)
		core.ErrorHandler(err3, "AllStrategyInitiativeCommunities.StrategyCommunityOrEmptyByID", "Could not StrategyCommunityOrEmptyByID")
		if found {
			if stratComm.ChannelCreated == 1 {
				op = append(op, each)
			}
		} else {
			log.Printf("Couldn't find the respective StrategyCommunity for id=%s", each.ID)
		}
	}
	return op
}

// UserStrategyInitiativeCommunities returns initiative communities associated with a user
// If the user is part of the strategy community, we return all initiative communities
// Else we return only those initiative communities that the user is a part of
func UserStrategyInitiativeCommunities(userID,
	communityUsersTable, communityUsersUserCommunityIndex, communityUsersUserIndex string,
	initiativeCommunitiesTable, initiativeCommunitiesPlatformIndex, strategyCommunitiesTable string,
	platformID models.PlatformID) []StrategyInitiativeCommunity {
	allInitiativeCommunities := AllStrategyInitiativeCommunities(platformID, initiativeCommunitiesTable,
		initiativeCommunitiesPlatformIndex, strategyCommunitiesTable)
	var op []StrategyInitiativeCommunity
	if community.IsUserInCommunity(userID, communityUsersTable, communityUsersUserCommunityIndex, community.Strategy) {
		// When user is in strategy community, list all the initiatives
		op = allInitiativeCommunities
	} else {
		// When the user is not part of strategy community, list only the initiative communities that the user is a part of
		// Getting the list of capability and initiative communities that user is a part of
		capComms, initComms := UserCapabilityInitiativeCommunities(userID, communityUsersTable, communityUsersUserIndex)
		var capCommIDs []string
		var initCommIDs []string
		for _, each := range capComms {
			ids := strings.Split(each.CommunityId, ":")
			capCommIDs = append(capCommIDs, ids[1])
		}
		for _, each := range initComms {
			ids := strings.Split(each.CommunityId, ":")
			initCommIDs = append(initCommIDs, ids[1])
		}
		for _, each := range allInitiativeCommunities {
			if core.ListContainsString(initCommIDs, each.ID) || core.ListContainsString(capCommIDs, each.CapabilityCommunityID) {
				op = append(op, each)
			}
		}
	}
	return op
}

func dynString(str string) *dynamodb.AttributeValue {
	attr := dynamodb.AttributeValue{S: aws.String(str)}
	return &attr
}

func getByIDAndPlatformIDUnsafe(table string, ID string, platformID models.PlatformID, result interface{}) {
	err2 := common.DeprecatedGetGlobalDns().Dynamo.GetItemFromTable(table, map[string]*dynamodb.AttributeValue{
		"id":          dynString(ID),
		"platform_id": dynString(string(platformID)),
	}, &result)
	core.ErrorHandler(err2, "getByIDAndPlatformIDUnsafe", fmt.Sprintf("Could not find %s in %s table", ID, table))
}

func CapabilityCommunityByID(platformID models.PlatformID, ID, table string) (res CapabilityCommunity) {
	getByIDAndPlatformIDUnsafe(table, ID, platformID, &res)
	if res.ID != ID {
		panic(fmt.Errorf("couldn't find CapabilityCommunityByID(ID=%s). Instead got ID=%s", ID, res.ID))
	}
	return
}

func StrategyObjectiveByID(platformID models.PlatformID, ID, table string) (res models.StrategyObjective) {
	getByIDAndPlatformIDUnsafe(table, ID, platformID, &res)
	if res.ID != ID {
		panic(fmt.Errorf("couldn't find StrategyObjectiveByID(ID=%s). Instead got ID=%s", ID, res.ID))
	}
	return
}

func StrategyInitiativeByID(platformID models.PlatformID, ID, table string) (res models.StrategyInitiative) {
	getByIDAndPlatformIDUnsafe(table, ID, platformID, &res)
	if res.ID != ID {
		panic(fmt.Errorf("couldn't find StrategyInitiativeByID(ID=%s). Instead got ID=%s", ID, res.ID))
	}
	return
}

func InitiativeCommunityByID(platformID models.PlatformID, ID, table string) (res StrategyInitiativeCommunity) {
	getByIDAndPlatformIDUnsafe(table, ID, platformID, &res)
	if res.ID != ID {
		panic(fmt.Errorf("couldn't find InitiativeCommunityByID(ID=%s). Instead got ID=%s", ID, res.ID))
	}
	return
}
// StrategyVision returns vision for platform ID or nil if absent
func StrategyVision(platformID models.PlatformID, visionTable string) (res *models.VisionMission) {
	log.Println("### In StrategyVision: platformID: " + platformID)
	// Query for the vision
	params := map[string]*dynamodb.AttributeValue{
		"platform_id": dynString(string(platformID)),
	}
	var vision models.VisionMission
	found, err2 := common.DeprecatedGetGlobalDns().Dynamo.GetItemOrEmptyFromTable(visionTable, params, &vision)
	core.ErrorHandler(err2, "StrategyVision", fmt.Sprintf("Could not find vision for platformID=%s in %s table", platformID, visionTable))
	if found {
		res = &vision
	} else {
		res =  nil
	}
	return
}

// String representation of new and old field values
func NewAndOld(new, old string) string {
	op := fmt.Sprintf("`New` - %s \n `Old` - %s", new, old)
	if new == old {
		op = new
	}
	return op
}

func communityEditMessage(typ community.AdaptiveCommunity, editStatus string) ui.RichText {
	return ui.RichText(fmt.Sprintf("This is the %s community you %s", typ, editStatus))
}

// CapabilityObjectivesDueInAWeek returns open capability objectives that exist for a user that are due in exactly 7 days
func CapabilityObjectivesDueInAWeek(userID string, ipDate business_time.Date, userObjectivesTable,
	userObjectivesUserIndex string) []models.UserObjective {
	return objectives.ObjectivesDueInDuration(userID, userObjectivesTable, userObjectivesUserIndex,
		models.StrategyDevelopmentObjective, models.ObjectiveStrategyObjectiveAlignment, ipDate, 7)
}

// CapabilityObjectivesDueInAMonth returns any open capability objectives that exist for a user that are due in exactly in 30 days
func CapabilityObjectivesDueInAMonth(userID string, ipDate business_time.Date, userObjectivesTable,
	userObjectivesUserIndex string) []models.UserObjective {
	return objectives.ObjectivesDueInDuration(userID, userObjectivesTable, userObjectivesUserIndex,
		models.StrategyDevelopmentObjective, models.ObjectiveStrategyObjectiveAlignment, ipDate, 30)
}

// CapabilityObjectivesDueInAQuarter returns any open capability objectives that exist for the user that are due in exactly in 90 days
func CapabilityObjectivesDueInAQuarter(userID string, ipDate business_time.Date, userObjectivesTable,
	userObjectivesUserIndex string) []models.UserObjective {
	return objectives.ObjectivesDueInDuration(userID, userObjectivesTable, userObjectivesUserIndex,
		models.StrategyDevelopmentObjective, models.ObjectiveStrategyObjectiveAlignment, ipDate, 90)
}

// Initiatives
// InitiativesDueInAWeek returns any open capability objectives that exist for the user that are due in exactly 7 days
func InitiativesDueInAWeek(userID string, ipDate business_time.Date, userObjectivesTable,
	userObjectivesUserIndex string) []models.UserObjective {
	return objectives.ObjectivesDueInDuration(userID, userObjectivesTable, userObjectivesUserIndex,
		models.StrategyDevelopmentObjective, models.ObjectiveStrategyInitiativeAlignment, ipDate, 7)
}

// InitiativesDueInAMonth return any open initiatives that exist for the user that are due in exactly 30 days
func InitiativesDueInAMonth(userID string, ipDate business_time.Date, userObjectivesTable,
	userObjectivesUserIndex string) []models.UserObjective {
	return objectives.ObjectivesDueInDuration(userID, userObjectivesTable, userObjectivesUserIndex,
		models.StrategyDevelopmentObjective, models.ObjectiveStrategyInitiativeAlignment, ipDate, 30)
}

// InitiativesDueInAQuarter returns any open initiatives that exist for the user that are due in exactly 90 days
func InitiativesDueInAQuarter(userID string, ipDate business_time.Date, userObjectivesTable,
	userObjectivesUserIndex string) []models.UserObjective {
	return objectives.ObjectivesDueInDuration(userID, userObjectivesTable, userObjectivesUserIndex,
		models.StrategyDevelopmentObjective, models.ObjectiveStrategyInitiativeAlignment, ipDate, 90)
}
