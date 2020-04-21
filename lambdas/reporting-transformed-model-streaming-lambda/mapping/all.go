package mapping

import (
	"github.com/jinzhu/gorm"
)

var allEntities = []interface{}{
	&DBClientConfig{},
	&DBUserFeedback{},
	&DBCommunityUser{},
	&DBCommunity{},
	&DBCompetency{},
	&DBEngagement{},
	&DBHoliday{},
	&DBInitiativeCommunity{},
	&DBObjectiveCommunity{},
	&DBObjective{},
	&DBPartnershipRejection{},
	&DBInitiative{},
	&DBUserObjectiveProgress{},
	&DBUserObjective{},
	&DBUser{},
	&DBVision{},
	&SlackTeam{},
}

// AutoMigrateAllEntities -
func AutoMigrateAllEntities(conn *gorm.DB) {
	conn.AutoMigrate(allEntities...)
}
