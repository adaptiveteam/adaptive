package lambda

import (
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	dialogFetcher "github.com/adaptiveteam/adaptive/dialog-fetcher"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	// utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
)

const (
	DateFormat      = core.ISODateLayout
	TimestampFormat = core.TimestampLayout
)

var (
	strategyEntityLabelMap = map[string]strategy.StrategyEntityDialog{
		fmt.Sprintf("%s:%s", strategy.StrategyCapabilityCommunityEntity, strategy.StrategyObjectiveEntity): {
			SourceLabel:      strategy.StrategyObjectiveCapabilityCommunityAssociationSourceLabel,
			TargetLabel:      strategy.StrategyObjectiveCapabilityCommunityAssociationTargetLabel,
			DescriptionLabel: strategy.StrategyObjectiveCapabilityCommunityAssociationDescriptionLabel,
		},
		fmt.Sprintf("%s:%s", strategy.StrategyObjectiveEntity, strategy.StrategyObjectiveEntity): {
			SourceLabel:      strategy.StrategyObjectivesAssociationSourceLabel,
			TargetLabel:      strategy.StrategyObjectivesAssociationTargetLabel,
			DescriptionLabel: strategy.StrategyObjectivesAssociationDescriptionLabel,
		},
		fmt.Sprintf("%s:%s", strategy.StrategyObjectiveEntity, strategy.StrategyInitiativeEntity): {
			SourceLabel:      strategy.StrategyObjectiveInitiativeAssociationSourceLabel,
			TargetLabel:      strategy.StrategyObjectiveInitiativeAssociationTargetLabel,
			DescriptionLabel: strategy.StrategyObjectiveInitiativeAssociationDescriptionLabel,
		},
		fmt.Sprintf("%s:%s", strategy.StrategyInitiativeCommunityEntity, strategy.StrategyInitiativeEntity): {
			SourceLabel:      strategy.StrategyInitiativeInitiativeCommunityAssociationSourceLabel,
			TargetLabel:      strategy.StrategyInitiativeInitiativeCommunityAssociationTargetLabel,
			DescriptionLabel: strategy.StrategyInitiativeInitiativeCommunityAssociationDescriptionLabel,
		},
	}

	strategyEntityTextMapping = map[strategy.StrategyEntityType]string{
		strategy.StrategyObjectiveEntity:           "strategy objectives",
		strategy.StrategyCapabilityCommunityEntity: "capability communities",
		strategy.StrategyInitiativeEntity:          "strategy initiatives",
		strategy.StrategyInitiativeCommunityEntity: "initiative communities",
	}
	schema              = models.SchemaForClientID(clientID)
	dialogTableName     = utils.NonEmptyEnv("DIALOG_TABLE")
	dialogFetcherDAO    = dialogFetcher.NewDAO(d, dialogTableName)
)

// allUsersInAnyStrategyCommunities should return users that belong to one of the communities
func allUsersInAnyStrategyCommunities(conn daosCommon.DynamoDBConnection) []models.KvPair {
	communityUsers := adaptiveCommunityUser.ReadByHashKeyPlatformIDUnsafe(conn.PlatformID)(conn)
	userIDsSet := getUserIDsSet(communityUsers)
	var users []models.KvPair
	// Get user options

	userProfiles := user.ReadAllUserProfiles(conn)
	for _, each := range userProfiles {
		if _, ok := userIDsSet[each.Id]; ok {
			users = append(users, models.KvPair{Key: each.DisplayName, Value: each.Id})

		}
	}
	return users
}

func getUserIDsSet(communityUsers []models.AdaptiveCommunityUser3) map[string]struct{} {
	userIDset := make(map[string]struct{})
	for _, usr := range communityUsers {
		userIDset[usr.UserID] = struct{}{}
	}
	return userIDset
}

func timeFormatChange(str string, oldFormat, newFormat core.AdaptiveDateLayout) string {
	t, err := oldFormat.ChangeLayout(str, newFormat)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse time %s using format %s", str, oldFormat))
	return t
}
