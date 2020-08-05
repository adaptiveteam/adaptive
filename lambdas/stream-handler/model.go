package streamhandler

import (
	"encoding/json"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/mapping"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	logger2 "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/jinzhu/gorm"
)

const ( 
	AdaptiveClientConfigTableRef          = "adaptive_client_config"
	AdaptiveUsersTableRef                 = "adaptive_users"
	AdaptiveCompetenciesTableRef          = "adaptive_value"
	AdaptiveVisionTableRef                = "vision"
	AdaptiveUserFeedbackTableRf           = "adaptive_user_feedback"
	AdaptiveHolidaysTableRef              = "ad_hoc_holidays"
	AdaptiveEngagementsTableRef           = "adaptive_users_engagements"
	AdaptiveUserObjectiveTableRef         = "user_objective"
	AdaptiveUserObjectiveProgressTableRef = "user_objectives_progress"
	AdaptiveCommunityTableRef             = "communities"
	AdaptiveObjectiveCommunityTableRef    = "capability_communities"
	AdaptiveCommunityUserTableRef         = "community_users"
	AdaptiveInitiativeCommunityTableRef   = "initiative_communities"
	AdaptivePartnershipRejectionTableRef  = "partnership_rejections" // does not exist
	AdaptiveInitiativeTableRef            = "strategy_initiatives"
	AdaptiveObjectiveTableRef             = "strategy_objectives"
	SlackTeamTableRef                     = "slack_team"
)

type InterfaceMapping func(interface{}, logger2.AdaptiveLogger) interface{}

type StreamHandling func(e2 model.StreamEntity, conn *gorm.DB, logger logger2.AdaptiveLogger)

// TODO: all these tables should have streaming enabled in .tf file:
// 	stream_enabled   = true
//	stream_view_type = var.dynamo_stream_view_type

var (
	TableMapping = map[string]mapping.DBEntity{
		AdaptiveClientConfigTableRef: mapping.DBClientConfig{},
		AdaptiveCommunityUserTableRef: mapping.DBCommunityUser{},
		AdaptiveCommunityTableRef: mapping.DBCommunity{},
		AdaptiveCompetenciesTableRef: mapping.DBCompetency{},
		AdaptiveEngagementsTableRef: mapping.DBEngagement{},
		AdaptiveHolidaysTableRef: mapping.DBHoliday{},
		AdaptiveInitiativeCommunityTableRef: mapping.DBInitiativeCommunity{},
		AdaptiveInitiativeTableRef: mapping.DBInitiative{},
		AdaptiveObjectiveCommunityTableRef: mapping.DBObjectiveCommunity{},
		AdaptiveObjectiveTableRef: mapping.DBObjective{},
		AdaptivePartnershipRejectionTableRef: mapping.DBPartnershipRejection{},
		SlackTeamTableRef: mapping.SlackTeam{},
		AdaptiveUserFeedbackTableRf: mapping.DBUserFeedback{},
		AdaptiveUserObjectiveProgressTableRef: mapping.DBUserObjectiveProgress{},
		AdaptiveUserObjectiveTableRef: mapping.DBUserObjective{},
		AdaptiveUsersTableRef: mapping.DBUser{},
		AdaptiveVisionTableRef: mapping.DBVision{},
	}
)

func TableRefKeys() (op []string) {
	for k := range TableMapping {
		op = append(op, k)
	}
	return
}

func StreamEntityHandler(
	entity model.StreamEntity,
	clientID string,
	conn1 *gorm.DB,
	logger logger2.AdaptiveLogger,
) (op model.StreamEntity) {
	op = entity
	tableSuffix := getTableSuffix(clientID, entity.TableName)

	logger.Infof("StreamEntityHandler(entity.TableName=%s)", entity.TableName)
	mappingFunction, found := TableMapping[tableSuffix]
	if found {
		oldJs, _ := json.Marshal(entity.OldEntity)
		op.OldEntity = mappingFunction.ParseUnsafe(oldJs, logger)
		newJs, _ := json.Marshal(entity.NewEntity)
		op.NewEntity = mappingFunction.ParseUnsafe(newJs, logger)

		mappingFunction.HandleStreamEntityUnsafe(op, conn1, logger)
	} else {
		logger.Warningf("Unhandled event from entity.TableName=%s", entity.TableName)
	}
	return
}

func getTableSuffix(clientID string, tableName string) string {
	return tableName[len(clientID) + 1 :]
}
