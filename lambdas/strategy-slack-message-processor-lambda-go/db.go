package lambda

import (
	"encoding/json"
	"fmt"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// StrategyInitiativeCommunityByID reads initiative community from table.
func StrategyInitiativeCommunityByID(id string, teamID models.TeamID) (result strategy.StrategyInitiativeCommunity) {
	entity := StrategyEntityById(id, teamID, strategyInitiativeCommunitiesTable)
	byt, _ := json.Marshal(entity)
	err := json.Unmarshal(byt, &result)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not decode interface to struct"))
	return
}

// StrategyCommunityByID reads community by ID (from `_strategy_communities` table)
// panics when not found.
// Deprecated. Use strategyCommunity.ReadOrEmptyUnsafe(id)(conn)
func StrategyCommunityByID(id string) (comm strategy.StrategyCommunity, found bool) {
	params := map[string]*dynamodb.AttributeValue{
		"id": daosCommon.DynS(id),
	}
	var err2 error
	found, err2 = d.GetItemOrEmptyFromTable(strategyCommunitiesTable, params, &comm)
	core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not query %s table", strategyCommunitiesTable))
	return
}

func SelectFromCapabilityCommunityJoinStrategyCommunityWhereChannelCreated(teamID models.TeamID) (out []strategy.CapabilityCommunity) {
	capComms := AllCapabilityCommunities(models.TeamID(teamID))
	conn := daosCommon.CreateConnectionFromEnv(teamID.ToPlatformID())
	for _, each := range capComms {
		comms := strategy.StrategyCommunityWithChannelByIDUnsafe(community.CapabilityPrefix, each.ID)(conn)
		if len(comms) >0 {
			out = append(out, each)
		}
	}
	return
}

func dynListString(list []string) *dynamodb.AttributeValue {
	return daosCommon.DynSS(list)
}

func dynString(str string) *dynamodb.AttributeValue {
	return daosCommon.DynS(str)
}

func dynInt(i int) *dynamodb.AttributeValue {
	return daosCommon.DynN(i)
}

func dynBool(b bool) *dynamodb.AttributeValue {
	return daosCommon.DynBOOL(b)
}

func idAndPlatformIDParams(id string, teamID models.TeamID) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue{
		"id":          dynString(id),
		"platform_id": dynString(teamID.ToString()),
	}
	return params
}

// Retrieve user objective by it's id
func userObjectiveByID(id string) models.UserObjective {
	return *objectives.UserObjectiveById(userObjectivesTable, id, dns)
}

func SetObjectiveField(uObj models.UserObjective, fieldName string, fieldValue interface{}) {
	var exprAttributes = map[string]*dynamodb.AttributeValue{}

	switch fieldValue.(type) {
	case int:
		exprAttributes = map[string]*dynamodb.AttributeValue{
			":f": dynInt(fieldValue.(int)),
		}
	case bool:
		exprAttributes = map[string]*dynamodb.AttributeValue{
			":f": dynBool(fieldValue.(bool)),
		}
	case string:
		exprAttributes = map[string]*dynamodb.AttributeValue{
			":f": dynString(fieldValue.(string)),
		}
	}

	key := map[string]*dynamodb.AttributeValue{
		"user_id": dynString(uObj.UserID),
		"id":      dynString(uObj.ID),
	}
	updateExpression := fmt.Sprintf("set %s = :f", fieldName)
	err := d.UpdateTableEntry(exprAttributes, key, updateExpression, userObjectivesTable)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not update field %s in %s table", fieldName, userObjectivesTable))
}

// userObjectiveProgressByID reads progress from db
// Set limit to -1 to retrieve all the updates
func userObjectiveProgressByID(id string, limit int) (ops []models.UserObjectiveProgress, err error) {
	// With scan forward to true, dynamo returns list in the ascending order of the range key
	err = d.QueryTableWithIndex(userObjectivesProgressTable, awsutils.DynamoIndexExpression{
		Condition: "id = :i",
		Attributes: map[string]interface{}{
			":i": id,
		},
	}, map[string]string{}, false, limit, &ops)
	return
}

// LatestProgressUpdateByObjectiveID retrieves the latest update, if exists, of by objective id
func LatestProgressUpdateByObjectiveID(id string) []models.UserObjectiveProgress {
	ops, err := userObjectiveProgressByID(id, 1)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not get progress for %s objective", id)
	}
	return ops
}

func getInitiativeCommunitiesForUserIDUnsafe(userID string, teamID models.TeamID) (initComms []strategy.StrategyInitiativeCommunity) {
	if isMemberInCommunity(userID, community.Strategy) {
		initComms = strategy.AllStrategyInitiativeCommunitiesWhereChannelExists(teamID)
	} else {
		initComms = StrategyInitiativeCommunitiesForUserID(userID, models.TeamID(teamID))
	}
	return
}
