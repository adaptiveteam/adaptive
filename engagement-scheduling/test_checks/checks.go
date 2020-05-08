package test_checks

import (
	adaptive_checks "github.com/adaptiveteam/adaptive/adaptive-checks"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/daos/common"
)

var ConstructTestProfile adaptive_checks.TypedProfileConstructor = func(conn common.DynamoDBConnection, userID string, date business_time.Date) adaptive_checks.TypedProfile {
	return adaptive_checks.SomeTrueAndSomeFalseTestProfile
}
