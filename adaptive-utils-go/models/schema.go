package models

import (
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/adaptiveteam/adaptive/daos/strategyObjective"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
)

// Schema describes tables and indices in DynamoDB
type Schema struct {
	Holidays HolidaysTableSchema
	AdaptiveValues AdaptiveValuesTableSchema
	ClientPlatformTokens ClientPlatformTokenTableSchema 
	AdaptiveUsers AdaptiveUsersTableSchema
	UserFeedback UserFeedbackTableSchema
	CommunityUsers CommunityUsersTableSchema
	AdaptiveCommunity AdaptiveCommunityTableSchema
}

// SchemaForClientID creates comprehensive schema for a given clientID
func SchemaForClientID(clientID string) Schema {
	return Schema{
		Holidays: HolidaysTableSchemaForClientID(clientID),
		AdaptiveValues: AdaptiveValuesTableSchemaForClientID(clientID),
		ClientPlatformTokens: ClientPlatformTokenTableSchemaForClientID(clientID),
		AdaptiveUsers: AdaptiveUsersTableSchemaForClientID(clientID),
		UserFeedback: UserFeedbackTableSchemaForClientID(clientID),
		CommunityUsers: CommunityUsersTableSchemaForClientID(clientID),
		AdaptiveCommunity: AdaptiveCommunityTableSchemaForClientID(clientID),
	}
}

var (
	DialogContentTableName                      = func(clientID string) string { return clientID + "_dialog_content" }
	StrategyObjectivesTableName                 = strategyObjective.TableName
	StrategyInitiativesTableName                = func(clientID string) string { return clientID + "_strategy_initiatives" }
	UserObjectivesTableName                     = userObjective.TableName
	UserObjectivesProgressTableName             = func(clientID string) string { return clientID + "_user_objectives_progress" }
	CommunityUsersTableName                     = adaptiveCommunityUser.TableName
	CommunitiesTableName                        = func(clientID string) string { return clientID + "_communities" }
	CompetenciesTableName                       = func(clientID string) string { return clientID + "_adaptive_value" }
	StrategyInitiativeCommunitiesTableName      = func(clientID string) string { return clientID + "_initiative_communities" }
	StrategyCommunityTableName                  = func(clientID string) string { return clientID + "_strategy_communities" }
	VisionTableName                             = func(clientID string) string { return clientID + "_vision" }
	CapabilityCommunitiesTableName              = func(clientID string) string { return clientID + "_capability_communities" }
	StrategyCommunitiesTableName                = func(clientID string) string { return clientID + "_strategy_communities" }
	AdaptiveUsersTableName                      = func(clientID string) string { return clientID + "_adaptive_users" }
)