package lambda

import (
	"encoding/json"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strconv"
)
// StrategyInitiativeCommunityByID reads initiative community from table.
func StrategyInitiativeCommunityByID(id string, platformID models.PlatformID) (result strategy.StrategyInitiativeCommunity) {
	entity := StrategyEntityById(id, platformID, strategyInitiativeCommunitiesTable)
	byt, _ := json.Marshal(entity)
	err := json.Unmarshal(byt, &result)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not decode interface to struct"))
	return
}

// StrategyCommunityByID reads community by ID (from `_strategy_communities` table)
// panics when not found.
func StrategyCommunityByID(id string) strategy.StrategyCommunity {
	params := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(id),
		},
	}
	var comm strategy.StrategyCommunity
	err := d.QueryTable(strategyCommunitiesTable, params, &comm)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not query %s table", strategyCommunitiesTable))
	return comm
}

func SelectFromCapabilityCommunityJoinStrategyCommunityWhereChannelCreated(platformID models.PlatformID) (out []strategy.CapabilityCommunity) {
	capComms := AllCapabilityCommunities(models.PlatformID(platformID))
	for _, each := range capComms {
		stratComm := StrategyCommunityByID(each.ID)
		if stratComm.ChannelCreated == 1 {
			out = append(out, each)
		}
	}
	return
}


func dynListString(list []string) *dynamodb.AttributeValue {
	attr := dynamodb.AttributeValue{SS: aws.StringSlice(list)}
	return &attr
}

func dynString(str string) *dynamodb.AttributeValue {
	attr := dynamodb.AttributeValue{S: aws.String(str)}
	return &attr
}

func dynInt(i int) *dynamodb.AttributeValue {
	attr := dynamodb.AttributeValue{N: aws.String(strconv.Itoa(i))}
	return &attr
}

func dynBool(b bool) *dynamodb.AttributeValue {
	attr := dynamodb.AttributeValue{BOOL: aws.Bool(b)}
	return &attr
}

func idAndPlatformIDParams(id string, platformID models.PlatformID) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue{
		"id":          dynString(id),
		"platform_id": dynString(string(platformID)),
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
		"id": dynString(uObj.ID),
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
