package lambda

import (
	"fmt"
	"strconv"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Retrieve user objective by it's id
func userObjectiveByID(id string) models.UserObjective {
	return *objectives.UserObjectiveById(userObjectivesTable, id, dns)
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
	ops, err2 := userObjectiveProgressByID(id, 1)
	if err2 != nil {
		logger.WithField("error", err2).Errorf("Could not get progress for %s objective", id)
	}
	return ops
}

// Get all the users associated with an accountability partner
// For a partner, this retrieves list of coachees
func UsersForPartner(partnerId string) []string {
	// Query all the objectives for which no partner is assigned
	var uObjs []models.UserObjective
	var users []string
	err2 := d.QueryTableWithIndex(userObjectivesTable, awsutils.DynamoIndexExpression{
		IndexName: string(userObjective.AccountabilityPartnerIndex),
		Condition: "accountability_partner = :ap",
		Attributes: map[string]interface{}{
			":ap": aws.String(partnerId),
		},
	}, map[string]string{}, true, -1, &uObjs)
	core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not query %s index on %s table", userObjective.AccountabilityPartnerIndex, userObjectivesTable))
	for _, each := range uObjs {
		users = append(users, each.UserID)
	}
	return core.Distinct(users)
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
	err2 := d.UpdateTableEntry(exprAttributes, key, updateExpression, userObjectivesTable)
	core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not update field %s in %s table", fieldName, userObjectivesTable))
}
