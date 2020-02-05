package issues

var (
	dialogContentTableName                      = func(clientID string) string { return clientID + "_dialog_content" }
	strategyObjectivesTableName                 = func(clientID string) string { return clientID + "_strategy_objectives" }
	strategyObjectivesPlatformIndex             = "StrategyObjectivesPlatformIndex"
	strategyInitiativesTableName                = func(clientID string) string { return clientID + "_strategy_initiatives" }
	strategyInitiativesPlatformIndex            = "StrategyInitiativesPlatformIndex"
	strategyInitiativesInitiativeCommunityIndex = "StrategyInitiativesInitiativeCommunityIndex"
	userObjectivesTableName                     = func(clientID string) string { return clientID + "_user_objective" }
	userObjectivesIDIndex                       = "IDIndex"
	userObjectivesUserIDIndex                   = "UserIDCompletedIndex"
	userObjectivesTypeIndex                     = "UserIDTypeIndex"
	userObjectivesProgressTableName             = func(clientID string) string { return clientID + "_user_objectives_progress" }
	communityUsersTableName                     = func(clientID string) string { return clientID + "_community_users" }
	communityUsersCommunityIndex                = "CommunityUsersCommunityIndex"
	communityUsersUserIndex                     = "CommunityUsersUserIndex"
	communitiesTableName                        = func(clientID string) string { return clientID + "_communities" }
	competenciesTableName                       = func(clientID string) string { return clientID + "_adaptive_value" }
	strategyInitiativeCommunitiesTableName      = func(clientID string) string { return clientID + "_initiative_communities" }
	strategyInitiativeCommunitiesPlatformIndex  = "PlatformIDIndex"
	strategyCommunityTableName                  = func(clientID string) string { return clientID + "_strategy_communities" }
	communityUsersUserCommunityIndex            = "UserIDCommunityIDIndex"
	visionTableName                             = func(clientID string) string { return clientID + "_vision" }
	capabilityCommunitiesTableName              = func(clientID string) string { return clientID + "_capability_communities" }
	capabilityCommunitiesPlatformIndex          = "CapabilityCommunitiesPlatformIndex"
	strategyCommunitiesTableName                = func(clientID string) string { return clientID + "_strategy_communities" }
	adaptiveUsersTableName                      = func(clientID string) string { return clientID + "_adaptive_users" }
	engagementTableName                         = func(clientID string) string { return clientID + "_adaptive_users_engagements" }

	objectiveCloseoutPath = ""
)
