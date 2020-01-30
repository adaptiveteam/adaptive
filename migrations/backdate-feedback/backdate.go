package main

import (
	"github.com/adaptiveteam/adaptive/daos/common"
	"log"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

func main() {
	defer core.RecoverAsLogError("main")
	config := readConfigFromEnvVars()

	feedback, err2 := readAllFeedbackForQuarterYear("2020","1", config)
	core.ErrorHandler(err2, config.namespace, "readAllFeedbackForQuarterYear")
	log.Printf("Backdating all feedback for platform id %s (1:2020 -> 4:2019)\n", config.platformID)
	res := []models.UserFeedback{}
	for _, f := range feedback {
		if f.PlatformID == config.platformID {
			f.QuarterYear = "4:2019"
			res = append(res, f)
		}
	}
	err3 := updateFeedback(res, config)
	core.ErrorHandler(err3, config.namespace, "writeFeedback")
}

// read feedback
func readAllFeedbackForQuarterYear(year, quarter string, config Config) (feedback []models.UserFeedback, err error) {
	err = config.d.QueryTableWithIndex(
		userFeedbackTableName(config.clientID), 
		awsutils.DynamoIndexExpression{
			IndexName: userFeedbackSourceQYIndex,
			Condition: "quarter_year = :qy",
			Attributes: map[string]interface{}{
				":qy": fmt.Sprintf("%d:%d", quarter, year),
			},
		}, 
		map[string]string{}, true, -1, &feedback)
	return
}

func updateFeedback(feedback []models.UserFeedback, config Config) (err error) {
	count := 0
	for _, f := range feedback {
		err = config.d.UpdateItemInTable(
			userFeedbackTableName(config.clientID),
			map[string]*dynamodb.AttributeValue {
				"qy": common.DynS(f.QuarterYear),
			},
			map[string]*dynamodb.AttributeValue {
				"id": common.DynS(f.ID),

			},
			"quarter_year = :qy",
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
