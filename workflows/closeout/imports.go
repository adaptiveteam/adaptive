package closeout

import (
	issues "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	common "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/workflows/exchange"
)

type Issue = issues.Issue
type DynamoDBConnection = common.DynamoDBConnection
type NewAndOldIssues = issues.NewAndOldIssues

// IssueIDKey - data key that will contain Issue ID
const IssueIDKey = exchange.IssueIDKey

// IssueTypeKey -
const IssueTypeKey = exchange.IssueTypeKey
