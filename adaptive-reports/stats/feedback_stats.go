package stats

import (
	"database/sql"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
)

const query = 
`SELECT sources/total*100 AS Given,
        targets/total*100 AS Received
FROM (SELECT COUNT(DISTINCT (user_feedback.source)) AS sources
      FROM user_feedback
      WHERE user_feedback.platform_id = ?
        AND year = ?
        AND quarter = ?) AS sources,
     (SELECT COUNT(DISTINCT (user_feedback.target)) AS targets
      FROM user_feedback
      WHERE user_feedback.platform_id = ?
        AND year = ?
        AND quarter = ?) AS targets,
     (SELECT COUNT(community_user.id) AS total
      FROM community_user
      WHERE community_user.community_id = 'user'
        AND community_user.platform_id = ?) AS total
`

type FeedbackStats struct {
	Given    float32
	Received float32
}
// QueryFeedbackStats queries RDS and returns feedback process statistics
func QueryFeedbackStats(
	teamID models.TeamID,
	year, quarter int,
) func (conn *sql.DB) (stats FeedbackStats, err error) {
	return func (conn *sql.DB) (stats FeedbackStats, err error) {
		platformID := teamID.ToPlatformID()
		var stmtOut *sql.Stmt
		stmtOut, err = conn.Prepare(query)
		defer stmtOut.Close()
		if err == nil {
			queryResult := stmtOut.QueryRow( 
				platformID, year, quarter, 
				platformID, year, quarter, 
				platformID,
			)
			if queryResult != nil {
				err = queryResult.Scan(&stats.Given, &stats.Received)
			}
		}
		return
	}
}
