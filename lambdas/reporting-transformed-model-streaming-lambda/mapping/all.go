package mapping

import (
	"github.com/jinzhu/gorm"
)

var allEntities = []interface{}{
	&DBClientConfig{},
	&DBCoachingFeedback{},
	&DBCommunityUser{},
	&DBCommunity{},
	&DBCompetency{},
	&DBEngagement{},
	&DBHoliday{},
	&DBInitiativeCommunity{},
	&DBObjectiveCommunity{},
	&DBObjective{},
	&DBPartnershipRejection{},
	&DBStrategyInitiative{},
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
