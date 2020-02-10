package workflow

import (
	"log"
	"time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/daos/postponedEvent"
	"github.com/adaptiveteam/adaptive/daos/common"
	
)

// PostponeEventHandler is a default postpone handler that will save to a database.
func PostponeEventHandler(conn common.DynamoDBConnection) func (platformID models.PlatformID, postponeEvent PostponeEventForAnotherUser) (err error) {
	dao := postponedEvent.NewDAO(conn.Dynamo, "PostponeEventHandler", conn.ClientID)
	return func (platformID models.PlatformID, postponeEvent PostponeEventForAnotherUser) (err error) {
		evt := postponedEvent.PostponedEvent{
			ID: core.Uuid(),
			PlatformID: platformID,
			UserID: postponeEvent.UserID,
			ActionPath: postponeEvent.ActionPath.Encode(),
			ValidThrough: core.TimestampLayout.Format(postponeEvent.ValidThrough),			
		}
		return dao.CreateOrUpdate(evt)
	}
}

// GetActionPathsForUserID returns all valid action paths for the current user.
func GetActionPathsForUserID(userID string) func(conn common.DynamoDBConnection) (actionPaths []models.ActionPath) {
	return func(conn common.DynamoDBConnection) (actionPaths []models.ActionPath) {
		dao := postponedEvent.NewDAO(conn.Dynamo, "GetActionPathsForUserID", conn.ClientID)
		events, err := dao.ReadByUserID(userID)
		now := time.Now()
		if err == nil {
			for _, e := range events {
				var validThrough time.Time
				validThrough, err = core.TimestampLayout.Parse(e.ValidThrough)
				if err != nil {
					return
				}
				if validThrough.After(now) {
					actionPaths = append(actionPaths, models.ParseActionPath(e.ActionPath))
				} else {
					log.Printf("Eliminating elapsed action for user %s, path=%s", userID, e.ActionPath)
					err = dao.Delete(e.ID)
					if err != nil {
						return
					}
				}
			}
		}
		return
	}
}

// ForeachActionPathForUserID runs the given function for all valid action paths of the current user.
func ForeachActionPathForUserID(userID string, f func(models.ActionPath, common.DynamoDBConnection)error) func(conn common.DynamoDBConnection) (count int, err error) {
	return func(conn common.DynamoDBConnection) (count int, err error) {
		dao := postponedEvent.NewDAO(conn.Dynamo, "ForeachActionPathForUserID", conn.ClientID)
		var events []postponedEvent.PostponedEvent
		events, err = dao.ReadByUserID(userID)
		now := time.Now()
		if err == nil {
			count = len(events)
			for _, e := range events {
				var validThrough time.Time
				validThrough, err = core.TimestampLayout.Parse(e.ValidThrough)
				if err != nil {
					return
				}
				if validThrough.After(now) {
					ap := models.ParseActionPath(e.ActionPath)
					err = f(ap, conn)
					if err != nil {
						return
					}
				} else {
					log.Printf("Eliminating elapsed action for user %s, path=%s", userID, e.ActionPath)
				}
				err = dao.Delete(e.ID)
				if err != nil {
					return
				}
			}
		}
		return
	}
}
