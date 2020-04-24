package issues

import (
	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"github.com/adaptiveteam/adaptive/daos/userObjectiveProgress"
	"github.com/pkg/errors"
)

// IssueProgressReadAll reads at most `limit` progress elements in descending order.
// Set limit to -1 to retrieve all the updates
func IssueProgressReadAll(issueID string, limit int) func(conn DynamoDBConnection) (res []userObjectiveProgress.UserObjectiveProgress, err error) {
	return func(conn DynamoDBConnection) (res []userObjectiveProgress.UserObjectiveProgress, err error) {
		// With scan forward to true, dynamo returns list in the ascending order of the range key
		scanForward := false
		err = conn.Dynamo.QueryTableWithIndex(
			models.UserObjectivesProgressTableName(conn.ClientID),
			awsutils.DynamoIndexExpression{
				Condition: "id = :i",
				Attributes: map[string]interface{}{
					":i": issueID,
				},
			}, map[string]string{}, scanForward, limit, &res)
		err = errors.Wrapf(err, "IssueProgressDynamoDBConnection) ReadAll(issueID=%s)", issueID)
		return
	}
}

func IssueProgressRead(issueProgressID IssueProgressID) func(conn DynamoDBConnection) (res userObjectiveProgress.UserObjectiveProgress, err error) {
	return func(conn DynamoDBConnection) (res userObjectiveProgress.UserObjectiveProgress, err error) {
		var ops []userObjectiveProgress.UserObjectiveProgress
		ops, err = userObjectiveProgress.ReadOrEmpty(issueProgressID.IssueID, issueProgressID.Date)(conn)
		if err == nil {
			if len(ops) > 0 {
				res = ops[0]
			} else {
				err = errors.New("UserObjectiveProgress " + issueProgressID.IssueID + " d: " + issueProgressID.Date + " not found")
			}
		}
		err = errors.Wrapf(err, "IssueProgressDynamoDBConnection) Read(issueProgressID=%s)", issueProgressID)
		return
	}
}
