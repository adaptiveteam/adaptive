package main

import (
	"github.com/adaptiveteam/adaptive/daos/common"
	"log"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

func main() {
	defer core.RecoverAsLogError("main")
	config := readConfigFromEnvVars()
	oldValue := "1:2020"
	newValue := "4:2019"
	feedback, err2 := readAllFeedbackForQuarterYear(oldValue, config)
	core.ErrorHandler(err2, config.namespace, "readAllFeedbackForQuarterYear")
	log.Printf("Backdating all feedback for platform id %s (%s -> %s)\n", config.platformID, oldValue, newValue)
	res := []models.UserFeedback{}
	for _, f := range feedback {
		if f.PlatformID == config.platformID {
			f.QuarterYear = newValue
			res = append(res, f)
		}
	}
	err3 := updateFeedback(res, config)
	core.ErrorHandler(err3, config.namespace, "writeFeedback")
}

// read feedback
func readAllFeedbackForQuarterYear(quarterYear string, config Config) (feedback []models.UserFeedback, err error) {
	err = config.d.QueryTableWithIndex(
		userFeedbackTableName(config.clientID), 
		awsutils.DynamoIndexExpression{
			IndexName: userFeedbackSourceQYIndex,
			Condition: "quarter_year = :qy",
			Attributes: map[string]interface{}{
				":qy": quarterYear, // fmt.Sprintf("%d:%d", quarter, year),
			},
		}, 
		map[string]string{}, true, -1, &feedback)
	return
}

func updateFeedback(feedback []models.UserFeedback, config Config) (err error) {
	count := 0
	for _, f := range feedback {
		key := map[string]*dynamodb.AttributeValue {
			"id": common.DynS(f.ID),
		}
		exprAttrs := map[string]*dynamodb.AttributeValue {
			":qy": common.DynS(f.QuarterYear),
		}
		err = config.d.UpdateItemInTable(
			userFeedbackTableName(config.clientID),
			key,
			"set quarter_year = :qy",
			exprAttrs,
		)
		log.Printf("Updated feedback %s", f.ID)
		count = count + 1
		if err != nil {
			return
		}
	}
	log.Printf("Updated %d feedback records", count)
	return
}
