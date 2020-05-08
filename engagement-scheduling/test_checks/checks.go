package test_checks

import (
	"github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/adaptive-checks"
)

var SomeTrueAndSomeFalseTestProfile = adaptive_checks.SomeTrueAndSomeFalseTestProfile

var ConstructTrueProfile adaptive_checks.TypedProfileConstructor = func (conn common.DynamoDBConnection, userID string, date business_time.Date) adaptive_checks.TypedProfile {
	return SomeTrueAndSomeFalseTestProfile
}
