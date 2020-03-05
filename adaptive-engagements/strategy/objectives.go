package strategy

import (
	"fmt"
	"strings"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

// QueryCommunityUserIndex gets Adaptive communities that the user is a part of
func QueryCommunityUserIndex(userId, table, index string) []models.AdaptiveCommunityUser2 {
	var rels []models.AdaptiveCommunityUser2
	err := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(table, awsutils.DynamoIndexExpression{
		IndexName: index,
		Condition: "user_id = :u",
		Attributes: map[string]interface{}{
			":u": userId,
		},
	}, map[string]string{}, true, -1, &rels)
	core.ErrorHandler(err, common.DeprecatedGetGlobalDns().Namespace, fmt.Sprintf("Could not query %s table on %s index", table, index))
	return rels
}

type DynamoTableWithIndex struct {
	Table string `json:"table"`
	Index string `json:"index"`
}

// StrategyCommunitiesDAOReadByPlatformID is a copy of daos/StrategyCommunity.go/DAO/ReadByPlatformID
func StrategyCommunitiesDAOReadByPlatformID(teamID models.TeamID, strategyCommunityTableName string) (out []StrategyCommunity, err error) {
	var instances []StrategyCommunity
	err = common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(strategyCommunityTableName,
		awsutils.DynamoIndexExpression{
			IndexName: "PlatformIDIndex",
			Condition: "platform_id = :a0",
			Attributes: map[string]interface{}{
				":a0": teamID.ToString(),
			},
		}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}

func UserCapabilityInitiativeCommunities(userID string, communityUsersTable, communityUsersUserIndex string) (
	[]models.AdaptiveCommunityUser2, []models.AdaptiveCommunityUser2) {
	var initComms []models.AdaptiveCommunityUser2
	var capComms []models.AdaptiveCommunityUser2
	// Get list of initiative communities for the user
	// Get initiates associated with those
	// Get objectives associated with initiated community
	commUsers := QueryCommunityUserIndex(userID, communityUsersTable, communityUsersUserIndex)
	for _, each := range commUsers {
		ids := strings.Split(each.CommunityId, ":")
		if len(ids) == 2 {
			switch ids[0] {
			case string(community.Initiative):
				initComms = append(initComms, each)
			case string(community.Capability):
				capComms = append(capComms, each)
			}
		}
	}

	return capComms, initComms
}

// UserInitiativeCommunityInitiatives lists all initiatives that are associated with initiative communities that the user is a part of
func UserInitiativeCommunityInitiatives(userID string, initiativesTableName, initiativesInitiativeCommunityIDIndex string,
	communityUsersTable, communityUsersUserIndex string) []models.StrategyInitiative {
	var op []models.StrategyInitiative
	// Get all initiative communities for a user
	_, initComms := UserCapabilityInitiativeCommunities(userID, communityUsersTable, communityUsersUserIndex)
	for _, each := range initComms {
		ids := strings.Split(each.CommunityId, ":")
		var inits []models.StrategyInitiative
		err := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(initiativesTableName, awsutils.DynamoIndexExpression{
			IndexName: initiativesInitiativeCommunityIDIndex,
			Condition: "initiative_community_id = :cc",
			Attributes: map[string]interface{}{
				":cc": ids[1],
			},
		}, map[string]string{}, true, -1, &inits)
		core.ErrorHandler(err, common.DeprecatedGetGlobalDns().Namespace, fmt.Sprintf("Could not query %s index on %s table",
			initiativesInitiativeCommunityIDIndex, initiativesTableName))
		op = append(op, inits...)
	}
	return op
}

// UserCapabilityCommunityInitiatives returns all initiatives associated with capability objectives of a user
func UserCapabilityCommunityInitiatives(userID string,
	strategyObjectivesTable, strategyObjectivesPlatformIndex, initiativesTable, initiativesInitiativeCommunityIDIndex,
	userObjectivesTable string,
	communityUsersTable, communityUsersUserCommunityIndex, communityUsersUserIndex string) []models.StrategyInitiative {
	var op []models.StrategyInitiative
	// Initiatives are associated with capability objectives
	capCommObjs := UserCommunityObjectives(userID, strategyObjectivesTable, strategyObjectivesPlatformIndex, userObjectivesTable,
		communityUsersTable, communityUsersUserCommunityIndex)
	var capObjIDs []string
	for _, each := range capCommObjs {
		capObjIDs = append(capObjIDs, each.ID)
	}
	// Get all capability communities for a user
	_, initComms := UserCapabilityInitiativeCommunities(userID, communityUsersTable, communityUsersUserIndex)
	for _, each := range initComms {
		ids := strings.Split(each.CommunityId, ":")

		var inits []models.StrategyInitiative
		err := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(initiativesTable, awsutils.DynamoIndexExpression{
			IndexName: initiativesInitiativeCommunityIDIndex,
			Condition: "initiative_community_id = :cc",
			Attributes: map[string]interface{}{
				":cc": ids[1],
			},
		}, map[string]string{}, true, -1, &inits)
		core.ErrorHandler(err, common.DeprecatedGetGlobalDns().Namespace, fmt.Sprintf("Could not query %s index on %s table",
			initiativesInitiativeCommunityIDIndex, initiativesTable))

		for _, ieach := range inits {
			if core.ListContainsString(capObjIDs, ieach.CapabilityObjective) {
				op = append(op, ieach)
			}
		}
	}
	return op
}

// UserCommunityObjectives lists all capability objectives that are associated with capability communities that the user is a part of
func UserCommunityObjectives(userID string, strategyObjectivesTable, strategyObjectivesPlatformIndex string,
	userObjectivesTable string,
	communityUsersTable, communityUsersUserIndex string) []models.StrategyObjective {
	var op []models.StrategyObjective
	teamID := UserIDToTeamID(userDAO())(userID)

	// Get all capability communities for a user
	capObjs := AllOpenStrategyObjectives(teamID, strategyObjectivesTable, strategyObjectivesPlatformIndex,
		userObjectivesTable)
	capComms, _ := UserCapabilityInitiativeCommunities(userID, communityUsersTable, communityUsersUserIndex)
	var added []string
	// Second, showing objectives
	for _, each := range capComms {
		ids := strings.Split(each.CommunityId, ":")
		for _, eachObj := range capObjs {
			if core.ListContainsString(eachObj.CapabilityCommunityIDs, ids[1]) {
				if !core.ListContainsString(added, eachObj.ID) {
					op = append(op, eachObj)
					// Appending only unique objectives
					added = append(added, eachObj.ID)
				}
			}
		}
	}
	return op
}

// UserCommunityInitiativesObjectives lists out all the objectives that are associated with capability communities and
// initiative communities that user is a part of
func UserCommunityInitiativesObjectives(userID string, strategyObjectivesTable, strategyObjectivesPlatformIndex string,
	userObjectivesTable string,
	initiativesTableName, initiativesInitiativeCommunityIDIndex string,
	communityUsersTable, communityUsersUserIndex string) []models.KvPair {
	var res = []models.KvPair{
		{
			Key:   "None",
			Value: "none",
		},
	}
	inits := UserInitiativeCommunityInitiatives(userID, initiativesTableName, initiativesInitiativeCommunityIDIndex, communityUsersTable, communityUsersUserIndex)
	for _, each := range inits {
		res = append(res, models.KvPair{Key: fmt.Sprintf("[%s] %s", strings.Title(string(community.Initiative)),
			each.Name), Value: fmt.Sprintf("%s:%s", community.Initiative, each.ID)})
	}

	objs := UserCommunityObjectives(userID, strategyObjectivesTable, strategyObjectivesPlatformIndex,
		userObjectivesTable,
		communityUsersTable, communityUsersUserIndex)
	for _, each := range objs {
		res = append(res, models.KvPair{Key: fmt.Sprintf("[%s] %s", strings.Title(string(community.Capability)),
			each.Name), Value: fmt.Sprintf("%s:%s", community.Capability, each.ID)})
	}

	return res
}

// UserAdvocacyObjectives gives a list of capability objectives that the user is an advocate for
func UserAdvocacyObjectives(userID, userObjectivesTable, userObjectivesTypeIndex string, completed int) []models.UserObjective {
	var op []models.UserObjective
	objs := objectives.UserObjectivesByType(userID, userObjectivesTable, userObjectivesTypeIndex,
		models.StrategyDevelopmentObjective, completed)
	for _, each := range objs {
		if each.StrategyAlignmentEntityType == models.ObjectiveStrategyObjectiveAlignment {
			op = append(op, each)
		}
	}
	return op
}

// UserAdvocacyInitiatives lists initiatives that the user is an advocate for
func UserAdvocacyInitiatives(userID, userObjectivesTable, userObjectivesTypeIndex string, completed int) []models.UserObjective {
	var op []models.UserObjective
	objs := objectives.UserObjectivesByType(userID, userObjectivesTable, userObjectivesTypeIndex,
		models.StrategyDevelopmentObjective, completed)
	for _, each := range objs {
		if each.StrategyAlignmentEntityType == models.ObjectiveStrategyInitiativeAlignment {
			op = append(op, each)
		}
	}
	return op
}

// UserCapabilityObjectivesWithNoProgressInAMonth retrieves all Capability Objectives for a user that haven't
// been updated in the last 30 days
func UserCapabilityObjectivesWithNoProgressInAMonth(userID string, ipDate business_time.Date,
	userObjectivesTable, userObjectivesUserIndex, userObjectivesProgressTable string, completed int) []models.UserObjective {
	aMonthBefore := ipDate.AddTime(0, 0, -30)
	fDay := aMonthBefore
	lDay := ipDate
	return objectives.StaleObjectivesInDuration(userID,
		userObjectivesTable, userObjectivesUserIndex, userObjectivesProgressTable,
		fDay, lDay, models.StrategyDevelopmentObjective, models.ObjectiveStrategyObjectiveAlignment)

}

// UserInitiativesWithNoProgressInAWeek retrieves all the Initiatives for a user that haven't been updated
// in the last 7 days
func UserInitiativesWithNoProgressInAWeek(userID string, ipDate business_time.Date, userObjectivesTable,
	userObjectivesUserIndex, userObjectivesProgressTable string, completed int) []models.UserObjective {
	aWeekBefore := ipDate.AddTime(0, 0, -7)
	fDay := aWeekBefore
	lDay := ipDate
	return objectives.StaleObjectivesInDuration(userID, userObjectivesTable,
		userObjectivesUserIndex, userObjectivesProgressTable,
		fDay, lDay,
		models.StrategyDevelopmentObjective, models.ObjectiveStrategyInitiativeAlignment)
}
