package streamhandler

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/mapping"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	logger2 "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/jinzhu/gorm"
	"strings"
)

const (
	AdaptiveClientConfigTableRef          = "adaptive_client_config"
	AdaptiveUsersTableRef                 = "adaptive_users"
	AdaptiveCompetenciesTableRef          = "adaptive_value"
	AdaptiveVisionTableRef                = "vision"
	AdaptiveCoachingFeedbackTableRf       = "adaptive_user_feedback"
	AdaptiveHolidaysTableRef              = "ad_hoc_holidays"
	AdaptiveEngagementsTableRef           = "adaptive_users_engagements"
	AdaptiveUserObjectiveTableRef         = "user_objective"
	AdaptiveUserObjectiveProgressTableRef = "user_objectives_progress"
	AdaptiveCommunityTableRef             = "communities"
	AdaptiveObjectiveCommunityTableRef    = "capability_communities"
	AdaptiveCommunityUserTableRef         = "community_users"
	AdaptiveInitiativeCommunityTableRef   = "initiative_communities"
	AdaptivePartnershipRejectionTableRef  = "partnership_rejections"
	AdaptiveInitiativeTableRef            = "strategy_initiatives"
	AdaptiveObjectiveTableRef             = "strategy_objectives"
	SlackTeamTableRef                     = "slack_team"
)

type InterfaceMapping func(interface{}, logger2.AdaptiveLogger) interface{}

type StreamHandling func(e2 model.StreamEntity, conn *gorm.DB, logger logger2.AdaptiveLogger)

type MappingFunction struct {
	TableMapped   string
	EntityMapper  InterfaceMapping
	StreamHandler StreamHandling
}

var (
	TableMapping = map[string]MappingFunction{
		AdaptiveClientConfigTableRef: {
			TableMapped:   "client_config", // done
			EntityMapper:  mapping.InterfaceToClientConfigUnsafe,
			StreamHandler: mapping.ClientConfigStreamEntityHandler,
		},
		AdaptiveUsersTableRef: {
			TableMapped:   "user", // done
			EntityMapper:  mapping.InterfaceToUserUnsafe,
			StreamHandler: mapping.UserStreamEntityHandler,
		},
		AdaptiveCompetenciesTableRef: {
			TableMapped:   "competency", // done
			EntityMapper:  mapping.InterfaceToCompetencyUnsafe,
			StreamHandler: mapping.CompetencyStreamEntityHandler,
		},
		AdaptiveVisionTableRef: {
			TableMapped:   "vision", // done
			EntityMapper:  mapping.InterfaceToVisionUnsafe,
			StreamHandler: mapping.VisionStreamEntityHandler,
		},
		AdaptiveCoachingFeedbackTableRf: {
			TableMapped:   "user_feedback",
			EntityMapper:  mapping.InterfaceToCoachingFeedbackUnsafe,
			StreamHandler: mapping.CoachingFeedbackStreamEntityHandler,
		},
		AdaptiveHolidaysTableRef: {
			TableMapped:   "holiday",
			EntityMapper:  mapping.InterfaceToHolidayUnsafe,
			StreamHandler: mapping.HolidayStreamEntityHandler,
		},
		AdaptiveEngagementsTableRef: {
			TableMapped:   "engagement", // done
			EntityMapper:  mapping.InterfaceToEngagementUnsafe,
			StreamHandler: mapping.EngagementStreamEntityHandler,
		},
		AdaptiveUserObjectiveTableRef: {
			TableMapped:   "user_objective", // done
			EntityMapper:  mapping.InterfaceToUserObjectiveUnsafe,
			StreamHandler: mapping.UserObjectiveStreamEntityHandler,
		},
		AdaptiveUserObjectiveProgressTableRef: {
			TableMapped:   "user_objective_progress", // done
			EntityMapper:  mapping.InterfaceToUserObjectiveProgressUnsafe,
			StreamHandler: mapping.UserObjectiveProgressStreamEntityHandler,
		},
		AdaptiveObjectiveCommunityTableRef: {
			TableMapped:   "objective_community", // done
			EntityMapper:  mapping.InterfaceToObjectiveCommunityUnsafe,
			StreamHandler: mapping.ObjectiveCommunityStreamEntityHandler,
		},
		AdaptiveCommunityTableRef: {
			TableMapped:   "community", //
			EntityMapper:  mapping.InterfaceToCommunityUnsafe,
			StreamHandler: mapping.CommunityStreamEntityHandler,
		},
		AdaptiveCommunityUserTableRef: {
			TableMapped:   "community_user", // done
			EntityMapper:  mapping.InterfaceToCommunityUserUnsafe,
			StreamHandler: mapping.CommunityUserStreamEntityHandler,
		},
		AdaptiveInitiativeCommunityTableRef: {
			TableMapped:   "initiative_community", // done
			EntityMapper:  mapping.InterfaceToInitiativeCommUnsafe,
			StreamHandler: mapping.InitiativeCommStreamEntityHandler,
		},
		AdaptivePartnershipRejectionTableRef: {
			TableMapped:   "partnership_rejection",
			EntityMapper:  mapping.InterfaceToPartnershipRejectionUnsafe,
			StreamHandler: mapping.PartnershipRejectionStreamEntityHandler,
		},
		AdaptiveInitiativeTableRef: {
			TableMapped:   "initiative", // done
			EntityMapper:  mapping.InterfaceToInitiativeUnsafe,
			StreamHandler: mapping.InitiativeStreamEntityHandler,
		},
		AdaptiveObjectiveTableRef: {
			TableMapped:   "objective", // done
			EntityMapper:  mapping.InterfaceToObjectiveUnsafe,
			StreamHandler: mapping.ObjectiveStreamEntityHandler,
		},
		SlackTeamTableRef: {
			TableMapped:   SlackTeamTableRef,
			EntityMapper:  mapping.InterfaceToSlackTeamUnsafe,
			StreamHandler: mapping.SlackTeamStreamEntityHandler,
		},
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
	op.EventType = entity.EventType
	tableName := entity.TableName

	found := false
	for _, each := range TableRefKeys() {
		if strings.Contains(tableName, fmt.Sprintf("%s_%s/", clientID, each)) {
			found = true
			mappingFunction := TableMapping[each]
			op.TableName = mappingFunction.TableMapped
			op.OldEntity = mappingFunction.EntityMapper(entity.OldEntity, logger)
			op.NewEntity = mappingFunction.EntityMapper(entity.NewEntity, logger)

			conn := conn1.Table(op.TableName)
			mappingFunction.StreamHandler(op, conn, logger)
		}
	}
	if !found {
		logger.Warningf("Unhandled event from tableName=%s", tableName)
	}
	return
}
