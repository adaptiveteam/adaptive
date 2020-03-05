package issues

import (
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
)

var (
	dialogContentTableName                      = models.DialogContentTableName
	strategyObjectivesTableName                 = models.StrategyObjectivesTableName
	strategyInitiativesTableName                = models.StrategyInitiativesTableName
	userObjectivesTableName                     = models.UserObjectivesTableName
	userObjectivesProgressTableName             = models.UserObjectivesProgressTableName
	communityUsersTableName                     = models.CommunityUsersTableName
	communitiesTableName                        = models.CommunitiesTableName
	competenciesTableName                       = models.CompetenciesTableName
	strategyInitiativeCommunitiesTableName      = models.StrategyInitiativeCommunitiesTableName
	strategyCommunityTableName                  = models.StrategyCommunityTableName
	visionTableName                             = models.VisionTableName
	capabilityCommunitiesTableName              = models.CapabilityCommunitiesTableName
	strategyCommunitiesTableName                = models.StrategyCommunitiesTableName
	adaptiveUsersTableName                      = models.AdaptiveUsersTableName

	objectiveCloseoutPath = ""
	strategyObjectivesPlatformIndex             = "PlatformIDIndex"
	strategyInitiativesPlatformIndex            = "PlatformIDIndex"
	strategyInitiativesInitiativeCommunityIndex = "InitiativeCommunityIDIndex"
	userObjectivesIDIndex                       = "IDIndex"
	userObjectivesUserIDIndex                   = "UserIDCompletedIndex"
	userObjectivesTypeIndex                     = "UserIDTypeIndex"
	communityUsersCommunityIndex                = string(adaptiveCommunityUser.PlatformIDCommunityIDIndex)
	communityUsersUserIndex                     = string(adaptiveCommunityUser.UserIDIndex)
	strategyInitiativeCommunitiesPlatformIndex  = "PlatformIDIndex"
	communityUsersUserCommunityIndex            = string(adaptiveCommunityUser.UserIDCommunityIDIndex)
	capabilityCommunitiesPlatformIndex          = "PlatformIDIndex"
)
