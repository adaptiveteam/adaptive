package stats

import (
	"database/sql"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
)

const query = 
`SELECT CONCAT(CAST(ROUND(sources/total*100,2) AS CHAR),'%') AS Given,
       CONCAT(CAST(ROUND(targets/total*100,2) AS CHAR),'%') AS Received
FROM (SELECT COUNT(DISTINCT (user_feedback.source)) AS sources
      FROM user_feedback
      WHERE user_feedback.platform_id = ?
        AND quarter = ?
        AND year = ?) AS sources,
     (SELECT COUNT(DISTINCT (user_feedback.target)) AS targets
      FROM user_feedback
      WHERE user_feedback.platform_id = ?
        AND quarter = ?
        AND year = ?) AS targets,
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
	quarter, year int,
) func (conn *sql.DB) (stats FeedbackStats, err error) {
	return func (conn *sql.DB) (stats FeedbackStats, err error) {
		platformID := teamID.ToPlatformID()
		var stmtOut *sql.Stmt
		stmtOut, err = conn.Prepare(query)
		defer stmtOut.Close()
		if err == nil {
			queryResult := stmtOut.QueryRow( 
				platformID, quarter, year, 
				platformID, quarter, year, 
				platformID,
			)
			if queryResult != nil {
				err = queryResult.Scan(&stats)
			}
		}
		return
	}
}
