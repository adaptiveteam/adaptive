package collaboration_report

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	fetch_dialog "github.com/adaptiveteam/adaptive/dialog-fetcher"
)

// BuildReport is an entry point for this project.
// Deprecated: Internally it uses global Dynamo DB
func BuildReport(
	// The last year of feedback received
	ReceivedBytes []byte,
	// The last year of feedback given
	GivenBytes []byte,
	// The users name (e.g., Chris Creel)
	UserName string,
	// The quarter for which this report was produced
	Quarter int,
	// The year for which this report was produced
	Year int,
	// Name and location for where to store the file.
	FileName string,
	logger logger.AdaptiveLogger,
	conn daosCommon.DynamoDBConnection,
) (tags map[string]string, err error) {
	dynamo := awsutils.NewDynamo(utils.NonEmptyEnv("AWS_REGION"), "", "dialog")
	dialogTableName := utils.NonEmptyEnv("DIALOG_TABLE")
	globalDao := fetch_dialog.NewDAO(dynamo, dialogTableName)
	tags, err = buildReport(
		ReceivedBytes,
		GivenBytes,
		UserName,
		Quarter,
		Year,
		FileName,
		globalDao,
		logger,
		conn,
	)
	return tags, err
}

// BuildReportWithCustomValues is an entry point for this project.
func BuildReportWithCustomValues(
	// The last year of feedback received
	ReceivedBytes []byte,
	// The last year of feedback given
	GivenBytes []byte,
	// The users name (e.g., Chris Creel)
	UserName string,
	// The quarter for which this report was produced
	Quarter int,
	// The year for which this report was produced
	Year int,
	// Name and location for where to store the file.
	FileName string,
	dialogDao fetch_dialog.DAO,
	logger logger.AdaptiveLogger,
	conn daosCommon.DynamoDBConnection,
) (tags map[string]string, err error) {
	tags, err = buildReport(
		ReceivedBytes,
		GivenBytes,
		UserName,
		Quarter,
		Year,
		FileName,
		dialogDao,
		logger,
		conn,
	)
	return tags, err
}
